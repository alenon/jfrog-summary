package commands

import (
	"encoding/json"
	"errors"
	tm "github.com/buger/goterm"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-cli-core/utils/coreutils"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ServerId        = "server-id"
	RefreshRate     = "refresh-rate"
	RecalculateRate = "recalculate-rate"
	MaximumResults  = "max-results"
)

type Summary struct {
	FileStoreSummary struct {
		StorageType      string `json:"storageType"`
		StorageDirectory string `json:"storageDirectory"`
		TotalSpace       string `json:"totalSpace"`
		UsedSpace        string `json:"usedSpace"`
		FreeSpace        string `json:"freeSpace"`
	} `json:"fileStoreSummary"`
	BinariesSummary struct {
		BinariesCount  string `json:"binariesCount"`
		BinariesSize   string `json:"binariesSize"`
		ArtifactsSize  string `json:"artifactsSize"`
		Optimization   string `json:"optimization"`
		ItemsCount     string `json:"itemsCount"`
		ArtifactsCount string `json:"artifactsCount"`
	} `json:"binariesSummary"`
	RepositoriesSummaryList []struct {
		RepoKey      string `json:"repoKey"`
		RepoType     string `json:"repoType"`
		FoldersCount int    `json:"foldersCount"`
		FilesCount   int    `json:"filesCount"`
		UsedSpace    string `json:"usedSpace"`
		ItemsCount   int    `json:"itemsCount"`
		PackageType  string `json:"packageType,omitempty"`
		Percentage   string `json:"percentage,omitempty"`
	} `json:"repositoriesSummaryList"`
}

type summaryConfiguration struct {
	refreshRate     int
	recalculateRate int
	maxResults      int
}

func GetStorageCommand() components.Command {
	return components.Command{
		Name:        "storage",
		Description: "Artifactory storage summary",
		Aliases:     []string{"st"},
		Flags:       getStorageFlags(),
		Action: func(c *components.Context) error {
			return storageCmd(c)
		},
	}
}

func getStorageFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:        ServerId,
			Description: "Artifactory server ID configured using the config command.",
		},
		components.StringFlag{
			Name:         RefreshRate,
			Description:  "Summary refresh rate in seconds",
			DefaultValue: "2",
		},
		components.StringFlag{
			Name:         RecalculateRate,
			Description:  "Storage summary recalculation rate in seconds. If 0 recalculation will not be triggered",
			DefaultValue: "0",
		},
		components.StringFlag{
			Name:         MaximumResults,
			Description:  "Maximal amount of shown results",
			DefaultValue: "10",
		},
	}
}

func storageCmd(c *components.Context) error {

	conf, err := prepareSummaryConf(c)
	if err != nil {
		return err
	}

	rtDetails, client, httpClientDetails, err := prepareHttpClient(c)
	if err != nil {
		return err
	}

	return fetchAndPresentSummary(conf, rtDetails, client, httpClientDetails)
}

func fetchAndPresentSummary(conf *summaryConfiguration, rtDetails *config.ArtifactoryDetails,
	client *rthttpclient.ArtifactoryHttpClient, httpClientDetails *httputils.HttpClientDetails) error {

	tm.Clear() // Clear current screen
	lastUpdated := time.Unix(0, 0)
	lastRecalculate := time.Unix(0, 0)
	for {
		if shouldRecalculate(conf, lastRecalculate) {
			go triggerRecalculate(rtDetails, client, httpClientDetails)
			lastRecalculate = time.Now()
		}
		if !shouldUpdateView(conf, lastUpdated) {
			continue
		}

		tm.MoveCursor(0, 0)
		_, _ = tm.Println("Last updated at:", time.Now().Format(time.RFC1123))
		if lastRecalculate.After(time.Unix(0, 0)) {
			_, _ = tm.Println("Last recalculated at:", lastRecalculate.Format(time.RFC1123))
		}

		_, _ = tm.Println("")
		err := showStorageSummary(conf, rtDetails, client, httpClientDetails)
		if err != nil {
			return err
		}
		tm.Flush() // Call it every time at the end of rendering
		lastUpdated = time.Now()
	}
}

func triggerRecalculate(rtDetails *config.ArtifactoryDetails, client *rthttpclient.ArtifactoryHttpClient,
	httpClientDetails *httputils.HttpClientDetails) {

	resp, _, err :=
		client.SendPost(rtDetails.GetUrl()+"api/storageinfo/calculate", nil, httpClientDetails)
	if err != nil {
		log.Error(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		log.Error(errors.New("Artifactory response: " + resp.Status))
	}
}

func showStorageSummary(conf *summaryConfiguration, rtDetails *config.ArtifactoryDetails, client *rthttpclient.ArtifactoryHttpClient,
	httpClientDetails *httputils.HttpClientDetails) error {

	storageSummary, err := fetchStorageSummary(rtDetails, client, httpClientDetails)
	if err != nil {
		return err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	configureTableView(table)

	totalsItem := storageSummary.RepositoriesSummaryList[len(storageSummary.RepositoriesSummaryList)-1]
	for i, row := range storageSummary.RepositoriesSummaryList {
		// limit on max results and don't print the row of totals
		if i >= conf.maxResults || row == totalsItem {
			table.SetFooter([]string{
				totalsItem.RepoKey,
				"-",
				"-",
				strconv.Itoa(totalsItem.FilesCount),
				totalsItem.UsedSpace,
				"-"})
			break
		}

		table.Append([]string{
			row.RepoKey,
			row.RepoType,
			row.PackageType,
			strconv.Itoa(row.FilesCount),
			row.UsedSpace,
			row.Percentage})
	}

	table.Render()
	_, _ = tm.Println(tableString.String())
	return nil
}

func configureTableView(table *tablewriter.Table) {
	table.SetHeader([]string{"Repository", "Type", "Package Type", " Files Count ", " Used Space ", " Percentage "})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.BgHiGreenColor},
		tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.BgHiGreenColor},
		tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.BgHiGreenColor},
		tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.BgHiGreenColor},
		tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.BgHiGreenColor},
		tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.BgHiGreenColor})

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor})

	table.SetColumnSeparator("│")
	table.SetCenterSeparator("┼")
	table.SetRowSeparator("─")
}

func fetchStorageSummary(rtDetails *config.ArtifactoryDetails, client *rthttpclient.ArtifactoryHttpClient,
	httpClientDetails *httputils.HttpClientDetails) (*Summary, error) {
	resp, respBody, _, err :=
		client.SendGet(rtDetails.GetUrl()+"api/storageinfo", false, httpClientDetails)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errorutils.CheckError(errors.New("Artifactory response: " + resp.Status))
	}

	var summaryStruct Summary
	err = json.Unmarshal(respBody, &summaryStruct)
	if err != nil {
		return nil, err
	}

	sort.Slice(summaryStruct.RepositoriesSummaryList, func(i, j int) bool {
		if summaryStruct.RepositoriesSummaryList[i].RepoKey == "TOTAL" {
			return false
		}
		return percentageSToI(summaryStruct.RepositoriesSummaryList[i].Percentage) >
			percentageSToI(summaryStruct.RepositoriesSummaryList[j].Percentage)
	})
	return &summaryStruct, nil
}

func prepareSummaryConf(c *components.Context) (*summaryConfiguration, error) {

	var summaryConfig = new(summaryConfiguration)
	refreshRate, err := strconv.Atoi(c.GetStringFlagValue(RefreshRate))
	if err != nil {
		return nil, errors.New("Illegal " + RefreshRate + " value. ")
	}
	summaryConfig.refreshRate = refreshRate

	recalculateRate, err := strconv.Atoi(c.GetStringFlagValue(RecalculateRate))
	if err != nil {
		return nil, errors.New("Illegal " + RecalculateRate + " value. ")
	}
	summaryConfig.recalculateRate = recalculateRate

	maxResults, err := strconv.Atoi(c.GetStringFlagValue(MaximumResults))
	if err != nil {
		return nil, errors.New("Illegal " + MaximumResults + " value. ")
	}
	summaryConfig.maxResults = maxResults
	return summaryConfig, nil
}

func prepareHttpClient(c *components.Context) (*config.ArtifactoryDetails,
	*rthttpclient.ArtifactoryHttpClient, *httputils.HttpClientDetails, error) {

	rtDetails, err := getRtDetails(c)
	if err != nil {
		return nil, nil, nil, err
	}

	auth, err := rtDetails.CreateArtAuthConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	securityDir, err := coreutils.GetJfrogSecurityDir()
	if err != nil {
		return nil, nil, nil, err
	}
	client, err := rthttpclient.ArtifactoryClientBuilder().
		SetCertificatesPath(securityDir).
		SetInsecureTls(rtDetails.InsecureTls).
		SetServiceDetails(&auth).
		Build()
	if err != nil {
		return nil, nil, nil, err
	}

	httpClientDetails := auth.CreateHttpClientDetails()
	return rtDetails, client, &httpClientDetails, nil
}

func getRtDetails(c *components.Context) (*config.ArtifactoryDetails, error) {

	serverId := c.GetStringFlagValue(ServerId)
	details, err := commands.GetConfig(serverId, false)
	if err != nil {
		return nil, err
	}
	if details.Url == "" {
		return nil, errors.New("no " + ServerId + " was found, or the " + ServerId + " has no url")
	}
	details.Url = clientutils.AddTrailingSlashIfNeeded(details.Url)
	err = config.CreateInitialRefreshableTokensIfNeeded(details)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func shouldRecalculate(conf *summaryConfiguration, lastRecalculate time.Time) bool {

	return conf.recalculateRate > 0 &&
		int(time.Since(lastRecalculate).Seconds()) >= conf.recalculateRate
}

func shouldUpdateView(conf *summaryConfiguration, lastUpdate time.Time) bool {

	return int(time.Since(lastUpdate).Seconds()) >= conf.refreshRate
}

func percentageSToI(percentage string) float64 {
	percentageFloatString := percentage[0 : len(percentage)-1]
	percentageFloat, err := strconv.ParseFloat(percentageFloatString, 64)
	if err != nil {
		log.Error("Failed parsing percentage value '" + percentage + "'")
	}
	return percentageFloat
}
