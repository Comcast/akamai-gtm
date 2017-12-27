package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/comcast/go-edgegrid/edgegrid"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// passed in via Makefile
var version string

func main() {
	app := cli.NewApp()
	app.Name = "akamai-gtm"
	app.Version = version
	app.Usage = "A CLI to Akamai GTM configuration"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Usage:  "Luna API Hostname",
			EnvVar: "AKAMAI_EDGEGRID_HOST",
		},
		cli.StringFlag{
			Name:   "client_token, ct",
			Usage:  "Luna API Client Token",
			EnvVar: "AKAMAI_EDGEGRID_CLIENT_TOKEN",
		},
		cli.StringFlag{
			Name:   "access_token, at",
			Usage:  "Luna API Access Token",
			EnvVar: "AKAMAI_EDGEGRID_ACCESS_TOKEN",
		},
		cli.StringFlag{
			Name:   "client_secret, s",
			Usage:  "Luna API Client Secret",
			EnvVar: "AKAMAI_EDGEGRID_CLIENT_SECRET",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "domains",
			Usage:       "domains",
			Description: "List all GTM Domains",
			Action:      domains,
		},
		{
			Name:        "domain",
			Usage:       "domain <domain.akadns.net>",
			Description: "View the details of a Domain",
			Action:      domain,
		},
		{
			Name:        "domain-create",
			Usage:       "domain-create --type <domainType> <domain.akadns.net>",
			Description: "Create a Domain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type",
					Usage: "The Domain type",
				},
			},
			Action: domainCreate,
		},
		{
			Name:        "domain-update",
			Usage:       "domain-update --json <DomainJSONFile>",
			Description: "Update a Domain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "The path to a JSON file",
				},
			},
			Action: domainUpdate,
		},
		{
			Name:        "data-centers",
			Usage:       "data-centers <domain.akadns.net>",
			Description: "List all DataCenters associated with a Domain",
			Action:      dataCenters,
		},
		{
			Name:        "data-centers-delete",
			Usage:       "data-centers-delete --id <dataCenterId> --id <dataCenterId> <domain.akadns.net>",
			Description: "Deletes specified DataCenters associated with a Domain",
			Flags: []cli.Flag{
				cli.IntSliceFlag{
					Name:  "id",
					Usage: "--id <dataCenterId> --id <dataCenterId>",
				},
			},
			Action: dataCentersDelete,
		},
		{
			Name:        "data-centers-delete-all",
			Usage:       "data-centers-delete-all <domain.akadns.net>",
			Description: "Deletes ALL DataCenters associated with a Domain",
			Action:      dataCentersDeleteAll,
		},
		{
			Name:        "data-center",
			Usage:       "data-center --id <dataCenterId> <domain.akadns.net>",
			Description: "View the details of a DataCenter associated with a Domain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The data center ID",
				},
			},
			Action: dataCenter,
		},
		{
			Name:        "data-center-create",
			Usage:       "data-center-create --json <DataCenterJSONFile> <domain.akadns.net>",
			Description: "Create a DataCenter associated with a Domain from data in a JSON file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "The path to a JSON file",
				},
			},
			Action: dataCenterCreate,
		},
		{
			Name:        "data-center-update",
			Usage:       "data-center-update --json <DataCenterJSONFile> <domain.akadns.net>",
			Description: "Update a DataCenter associated with a Domain from data in a JSON file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "The path to a JSON file",
				},
			},
			Action: dataCenterUpdate,
		},
		{
			Name:        "data-center-delete",
			Usage:       "data-center-delete --id <dataCenterId> <domain.akadns.net>",
			Description: "Delete a data center",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The data center ID",
				},
			},
			Action: dataCenterDelete,
		},
		{
			Name:        "properties",
			Usage:       "properties",
			Description: "View all Properties of a Domain",
			Action:      properties,
		},
		{
			Name:        "properties-delete",
			Usage:       "properties-delete --names <PropertyName>,<PropertyName> <domain.akadns.net>",
			Description: "Deletes specified Properties associated with a Domain",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "names",
					Usage: "A comma-separated list of names of Properties to delete",
				},
			},
			Action: propertiesDelete,
		},
		{
			Name:        "properties-delete-all",
			Usage:       "properties-delete-all <domain.akadns.net>",
			Description: "Deletes ALL Properties associated with a Domain",
			Action:      propertiesDeleteAll,
		},
		{
			Name:        "property",
			Usage:       "property --name <PropertyName> <domain.akadns.net>",
			Description: "View the details of a Property",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The Property name",
				},
			},
			Action: property,
		},
		{
			Name:        "property-create",
			Usage:       "property-create --json <PropertyJSONFile> <domain.akadns.net>",
			Description: "Create a Property",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "The path to a JSON file",
				},
			},
			Action: propertyCreate,
		},
		{
			Name:        "property-update",
			Usage:       "property-update --json <PropertyJSONFile> <domain.akadns.net>",
			Description: "Update a Property",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "The path to a JSON file",
				},
			},
			Action: propertyUpdate,
		},
		{
			Name:        "property-delete",
			Usage:       "property-delete --name <PropertyName> <domain.akadns.net>",
			Description: "Delete a Property",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The Property name",
				},
			},
			Action: propertyDelete,
		},
		{
			Name:        "traffic-targets",
			Usage:       "traffic-targets --name <PropertyName> <domain.akadns.net>",
			Description: "View traffic targets associated with a property",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The Property name",
				},
			},
			Action: trafficTargets,
		},
		{
			Name:        "liveness-tests",
			Usage:       "liveness-tests --name <PropertyName> <domain.akadns.net>",
			Description: "View liveness tests associated with a property",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The Property name",
				},
			},
			Action: livenessTests,
		},
		{
			Name:        "status",
			Usage:       "status <domain.akadns.net>",
			Description: "View the Status details for a Domain",
			Action:      status,
		},
	}
	app.Run(os.Args)
}

func domains(c *cli.Context) error {
	client := client(c)
	domains, err := client.Domains()
	if err != nil {
		return err
	}
	for _, domain := range domains {
		fmt.Printf("%s\n", domain.Name)
	}

	return nil
}

func domain(c *cli.Context) error {
	client := client(c)
	domain, err := client.Domain(c.Args().First())
	if err != nil {
		return err
	}
	dcs := []string{}
	for _, dc := range domain.Datacenters {
		dcs = append(dcs, dc.Nickname)
	}
	data := [][]string{
		[]string{"Name", domain.Name},
		[]string{"Type", domain.Type},
		[]string{"DataCenters", strings.Join(dcs, ", ")},
		[]string{"Status", domain.Status.Message},
		[]string{"Propagation Status", domain.Status.PropagationStatus},
		[]string{"Last Modified By", domain.LastModifiedBy},
		[]string{"Last Modified", domain.LastModified},
		[]string{"Modification Comments", domain.ModificationComments},
	}

	printBasicTable(data)

	return nil
}

func domainCreate(c *cli.Context) error {
	client := client(c)
	domainResp, err := client.DomainCreate(c.Args().First(), c.String("type"))
	if err != nil {
		return err
	}

	fmt.Printf("Created %s\n", domainResp.Domain.Name)

	return nil
}

func domainUpdate(c *cli.Context) error {
	client := client(c)
	domainSt := &edgegrid.Domain{}
	data, err := ioutil.ReadFile(c.String("json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, domainSt); err != nil {
		return err
	}

	domainResp, err := client.DomainUpdate(domainSt)
	if err != nil {
		return err
	}

	fmt.Printf("Updated domain: %s\n", domainResp.Domain.Name)

	return nil
}

func dataCenters(c *cli.Context) error {
	domain := c.Args().First()
	client := client(c)
	dcs, err := client.DataCenters(domain)
	if err != nil {
		return err
	}
	data := [][]string{}
	for _, dc := range dcs {
		data = append(data, []string{dc.Nickname, strconv.Itoa(dc.DataCenterID)})
	}

	if len(data) != 0 {
		printTableWithHeaders([]string{"Nickname", "DataCenter ID"}, data)
	} else {
		fmt.Printf("No data centers found for domain: %s\n", domain)
	}

	return nil
}

func dataCenter(c *cli.Context) error {
	dc, err := client(c).DataCenter(c.Args().First(), c.Int("id"))
	if err != nil {
		return err
	}
	data := [][]string{
		[]string{"Nickname", dc.Nickname},
		[]string{"DataCenterID", strconv.Itoa(dc.DataCenterID)},
		[]string{"City", dc.City},
		[]string{"CloneOf", strconv.Itoa(dc.CloneOf)},
		[]string{"Continent", dc.Continent},
		[]string{"Country", dc.Country},
		[]string{"Latitude", floatToStr(dc.Latitude)},
		[]string{"Longitude", floatToStr(dc.Longitude)},
		[]string{"StateOrProvince", dc.StateOrProvince},
		[]string{"Virtual", strconv.FormatBool(dc.Virtual)},
		[]string{"CloudServerTargeting", strconv.FormatBool(dc.CloudServerTargeting)},
	}

	printBasicTable(data)

	return nil
}

func dataCenterCreate(c *cli.Context) error {
	data := unmarshalDc(c)
	dc, err := client(c).DataCenterCreate(c.Args().First(), data)
	if err != nil {
		return err
	}

	fmt.Printf("Created %s\n", dc.DataCenter.Nickname)

	return nil
}

func dataCenterUpdate(c *cli.Context) error {
	data := unmarshalDc(c)
	dc, err := client(c).DataCenterUpdate(c.Args().First(), data)
	if err != nil {
		return err
	}

	fmt.Printf("Updated %s\n", dc.DataCenter.Nickname)

	return nil
}

func unmarshalDc(c *cli.Context) *edgegrid.DataCenter {
	dcSt := &edgegrid.DataCenter{}
	data, err := ioutil.ReadFile(c.String("json"))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, dcSt); err != nil {
		panic(err)
	}

	return dcSt
}

func dataCenterDelete(c *cli.Context) error {
	id := c.Int("id")
	err := client(c).DataCenterDelete(c.Args().First(), id)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted data center %d\n", id)

	return nil
}

func dataCentersDelete(c *cli.Context) error {
	ids := c.IntSlice("id")
	client := client(c)
	domainName := c.Args().First()

	for _, id := range ids {
		err := client.DataCenterDelete(domainName, id)
		if err != nil {
			fmt.Printf("Failed to delete DataCenter: %d\n", id)
			fmt.Printf("Error is: %v", err)

			return err
		}

		fmt.Printf("Deleted DataCenter: %d\n", id)
	}

	return nil
}

func dataCentersDeleteAll(c *cli.Context) error {
	myClient := client(c)
	domainName := c.Args().First()
	dcs, err := myClient.DataCenters(domainName)
	if err != nil {
		return err
	}
	for _, dc := range dcs {
		err := myClient.DataCenterDelete(domainName, dc.DataCenterID)
		if err != nil {
			fmt.Printf("Failed to delete DataCenter: %s\n", dc.Nickname)
			fmt.Printf("Error is: %v", err)

			return err
		}

		fmt.Printf("Deleted DC: %s (%d)\n", dc.Nickname, dc.DataCenterID)
	}

	return nil
}

// Properties is a Property slice
type Properties []Property

// Property is an Akamai GTM property
type Property struct {
	GtmProperty edgegrid.Property
	Product     string
}

func (props Properties) Len() int {
	return len(props)
}

func (props Properties) Less(i, j int) bool {
	return props[i].Product < props[j].Product
}

func (props Properties) Swap(i, j int) {
	props[i], props[j] = props[j], props[i]
}

func buildProps(props *edgegrid.Properties) Properties {
	created := Properties{}

	for _, prop := range props.Properties {
		splitName := strings.Split(prop.Name, ".")
		created = append(created, Property{
			GtmProperty: prop,
			Product:     splitName[len(splitName)-1],
		})
	}

	return created
}

func properties(c *cli.Context) error {
	domain := c.Args().First()
	ps, err := client(c).Properties(domain)
	if err != nil {
		return err
	}

	props := buildProps(ps)
	sort.Sort(props)

	s := [][]string{}
	for _, prop := range props {
		s = append(s, []string{
			prop.GtmProperty.Name,
			prop.GtmProperty.Type,
			strings.Join(targetIds(prop.GtmProperty.TrafficTargets), ", "),
		})
	}

	if len(s) != 0 {
		printTableWithHeaders([]string{"Name", "Type", "Traffic Targets"}, s)
	} else {
		fmt.Printf("No properties found for domain: %s\n", domain)
	}

	return nil
}

func property(c *cli.Context) error {
	prop, err := client(c).Property(c.Args().First(), c.String("name"))
	if err != nil {
		return err
	}

	printProp(prop)

	return nil
}

func propertyCreate(c *cli.Context) error {
	data := unmarshalProp(c)
	prop, err := client(c).PropertyCreate(c.Args().First(), data)
	if err != nil {
		return err
	}

	printProp(prop.Property)

	return nil
}

func propertyUpdate(c *cli.Context) error {
	data := unmarshalProp(c)
	prop, err := client(c).PropertyUpdate(c.Args().First(), data)
	if err != nil {
		return err
	}

	printProp(prop.Property)

	return nil
}

func unmarshalProp(c *cli.Context) *edgegrid.Property {
	propSt := &edgegrid.Property{}
	data, err := ioutil.ReadFile(c.String("json"))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, propSt); err != nil {
		panic(err)
	}

	return propSt
}

func propertyDelete(c *cli.Context) error {
	name := c.String("name")
	_, err := client(c).PropertyDelete(c.Args().First(), name)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted property %s\n", name)

	return nil
}

func trafficTargets(c *cli.Context) error {
	name := c.String("name")
	prop, err := client(c).Property(c.Args().First(), name)
	if err != nil {
		return err
	}

	for _, target := range prop.TrafficTargets {
		fmt.Printf("\nTraffic target\n")
		printTrafficTarget(target)
	}

	return nil
}

func printTrafficTarget(target edgegrid.TrafficTarget) error {
	data := [][]string{
		[]string{"Name", interfaceToStr(target.Name)},
		[]string{"DCId", strconv.Itoa(target.DataCenterID)},
		[]string{"Enabled", strconv.FormatBool(target.Enabled)},
		[]string{"HandoutCname", interfaceToStr(target.HandoutCname)},
		[]string{"Servers", strings.Join(target.Servers, ", ")},
		[]string{"Weight", floatToStr(target.Weight)},
	}

	printBasicTable(data)

	return nil
}

func livenessTests(c *cli.Context) error {
	prop := c.String("name")
	props, err := client(c).Property(c.Args().First(), prop)
	if err != nil {
		return err
	}

	for _, test := range props.LivenessTests {
		printLivenessTest(test)
	}

	return nil
}

func printLivenessTest(test edgegrid.LivenessTest) error {
	data := [][]string{
		[]string{"Name", test.Name},
		[]string{"HTTPError3xx", strconv.FormatBool(test.HTTPError3xx)},
		[]string{"HTTPError4xx", strconv.FormatBool(test.HTTPError4xx)},
		[]string{"HTTPError5xx", strconv.FormatBool(test.HTTPError5xx)},
		[]string{"TestInterval", strconv.FormatInt(test.TestInterval, 10)},
		[]string{"TestObject", test.TestObject},
		[]string{"TestObjectPort", strconv.FormatInt(test.TestObjectPort, 10)},
		[]string{"TestObjectProtocol", test.TestObjectProtocol},
		[]string{"TestObjectUsername", test.TestObjectUsername},
		[]string{"TestObjectPassword", test.TestObjectPassword},
		[]string{"TestTimeout", floatToStr(test.TestTimeout)},
		[]string{"DisableNonstandardPortWarning", strconv.FormatBool(test.DisableNonstandardPortWarning)},
		[]string{"RequestString", test.RequestString},
		[]string{"ResponseString", test.ResponseString},
		[]string{"SSLClientPrivateKey", test.SSLClientPrivateKey},
		[]string{"SSLCertificate", test.SSLCertificate},
		[]string{"HostHeader", test.HostHeader},
	}

	printBasicTable(data)

	return nil
}

func propertiesDelete(c *cli.Context) error {
	names := strings.Split(c.String("names"), ",")
	client := client(c)
	domain := c.Args().First()

	for _, name := range names {
		_, err := client.PropertyDelete(domain, strings.TrimSpace(name))
		if err != nil {
			fmt.Printf("Failed to delete Property: %s\n", name)
			fmt.Printf("Error is: %v", err)

			return err
		}

		fmt.Printf("Deleted Property: %s\n", name)
	}

	return nil
}

func propertiesDeleteAll(c *cli.Context) error {
	domain := c.Args().First()
	client := client(c)
	props, err := client.Properties(domain)
	if err != nil {
		return err
	}

	for _, prop := range props.Properties {
		_, err := client.PropertyDelete(domain, prop.Name)
		if err != nil {
			fmt.Printf("Failed to delete Property: %s\n", prop.Name)
			fmt.Printf("Error is: %v", err)
		} else {
			fmt.Printf("Deleted Property: %s\n", prop.Name)
		}
	}

	return nil
}

func client(c *cli.Context) *edgegrid.GTMClient {
	return edgegrid.GTMClientWithCreds(
		c.GlobalString("access_token"),
		c.GlobalString("client_token"),
		c.GlobalString("client_secret"),
		c.GlobalString("host"))
}

func targetIds(trafficTargets []edgegrid.TrafficTarget) []string {
	targets := []string{}

	for _, t := range trafficTargets {
		targets = append(targets, strconv.Itoa(t.DataCenterID))
	}

	return targets
}

func livenessTestNames(livenessTests []edgegrid.LivenessTest) []string {
	tests := []string{}

	for _, t := range livenessTests {
		tests = append(tests, t.Name)
	}

	return tests
}

func printProp(prop *edgegrid.Property) error {
	data := [][]string{
		[]string{"BackupCname", prop.BackupCname},
		[]string{"BackupIP", prop.BackupIP},
		[]string{"BalanceByDownloadScore", strconv.FormatBool(prop.BalanceByDownloadScore)},
		[]string{"Cname", prop.Cname},
		[]string{"Comments", prop.Comments},
		[]string{"DynamicTTL", strconv.Itoa(prop.DynamicTTL)},
		[]string{"FailbackDelay", strconv.Itoa(prop.FailbackDelay)},
		[]string{"FailoverDelay", strconv.Itoa(prop.FailoverDelay)},
		[]string{"HandoutMode", prop.HandoutMode},
		[]string{"HealthMax", floatToStr(prop.HealthMax)},
		[]string{"HealthMultiplier", floatToStr(prop.HealthMultiplier)},
		[]string{"HealthThreshold", floatToStr(prop.HealthThreshold)},
		[]string{"Ipv6", strconv.FormatBool(prop.Ipv6)},
		[]string{"LastModified", prop.LastModified},
		[]string{"LivenessTests", strings.Join(livenessTestNames(prop.LivenessTests), ", ")},
		[]string{"LoadImbalancePercentage", floatToStr(prop.LoadImbalancePercentage)},
		[]string{"MapName", interfaceToStr(prop.MapName)},
		[]string{"MaxUnreachablePenalty", interfaceToStr(prop.MaxUnreachablePenalty)},
		[]string{"MxRecords", maxRecsString(prop.MxRecords)},
		[]string{"Name", prop.Name},
		[]string{"ScoreAggregationType", prop.ScoreAggregationType},
		[]string{"StaticTTL", interfaceToStr(prop.StaticTTL)},
		[]string{"StickinessBonusConstant", interfaceToStr(prop.StickinessBonusConstant)},
		[]string{"StickinessBonusPercentage", interfaceToStr(prop.StickinessBonusPercentage)},
		[]string{"TrafficTargets", strings.Join(targetIds(prop.TrafficTargets), ", ")},
		[]string{"Type", prop.Type},
		[]string{"UnreachableThreshold", interfaceToStr(prop.UnreachableThreshold)},
		[]string{"UseComputedTargets", strconv.FormatBool(prop.UseComputedTargets)},
	}

	printBasicTable(data)

	return nil
}

func status(c *cli.Context) error {
	client := client(c)
	status, err := client.DomainStatus(c.Args().First())
	if err != nil {
		return err
	}
	data := [][]string{
		[]string{"PropagationStatus", status.PropagationStatus},
		[]string{"PassingValidation", strconv.FormatBool(status.PassingValidation)},
		[]string{"Message", status.Message},
		[]string{"ChangeID", status.ChangeID},
		[]string{"PropagationStatusDate", status.PropagationStatusDate},
	}

	printBasicTable(data)

	return nil
}

func printBasicTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(data)
	table.SetRowLine(true)
	table.Render()
}

func printTableWithHeaders(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(data)
	table.SetRowLine(true)
	table.Render()
}

func floatToStr(input float64) string {
	return strconv.FormatFloat(input, 'f', 6, 64)
}

func maxRecsString(records []interface{}) string {
	recs := []string{}

	for _, r := range records {
		recs = append(recs, r.(string))
	}

	return strings.Join(recs, ", ")
}

func interfaceToStr(inter interface{}) string {
	if str, ok := inter.(string); ok {
		return str
	}
	if num, ok := inter.(int); ok {
		return strconv.Itoa(num)
	}
	if float, ok := inter.(float64); ok {
		return floatToStr(float)
	}

	return ""
}
