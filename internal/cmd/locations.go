package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	threads "github.com/salmonumbrella/threads-go"
	"github.com/salmonumbrella/threads-go/internal/outfmt"
	"github.com/salmonumbrella/threads-go/internal/ui"
)

var locationsCmd = &cobra.Command{
	Use:     "locations",
	Aliases: []string{"location", "loc"},
	Short:   "Location search and details",
	Long: `Search for locations and retrieve location details.

Locations can be tagged on posts using the threads_location_tagging scope.
Use 'locations search' to find location IDs, then include them when creating posts.`,
}

var locationsSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for locations",
	Long: `Search for locations by name or coordinates.

Examples:
  threads locations search "Central Park"
  threads locations search --lat 40.7829 --lng -73.9654
  threads locations search "coffee" --lat 37.7749 --lng -122.4194`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLocationsSearch,
}

var locationsGetCmd = &cobra.Command{
	Use:   "get [location-id]",
	Short: "Get location details",
	Long: `Retrieve detailed information about a specific location.

Example:
  threads locations get 123456789`,
	Args: cobra.ExactArgs(1),
	RunE: runLocationsGet,
}

var (
	locLat float64
	locLng float64
)

func init() {
	locationsCmd.AddCommand(locationsSearchCmd)
	locationsCmd.AddCommand(locationsGetCmd)

	locationsSearchCmd.Flags().Float64Var(&locLat, "lat", 0, "Latitude for coordinate search")
	locationsSearchCmd.Flags().Float64Var(&locLng, "lng", 0, "Longitude for coordinate search")
}

func runLocationsSearch(cmd *cobra.Command, args []string) error {
	var query string
	if len(args) > 0 {
		query = args[0]
	}

	if query == "" && locLat == 0 && locLng == 0 {
		return fmt.Errorf("provide either a search query or --lat/--lng coordinates")
	}

	ctx := cmd.Context()
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	var latPtr, lngPtr *float64
	if locLat != 0 || locLng != 0 {
		latPtr = &locLat
		lngPtr = &locLng
	}

	result, err := client.SearchLocations(ctx, query, latPtr, lngPtr)
	if err != nil {
		return fmt.Errorf("location search failed: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(result, jqQuery)
	}

	if len(result.Data) == 0 {
		ui.Info("No locations found")
		return nil
	}

	ui.Success("Found %d location(s)", len(result.Data))
	fmt.Println()

	f := outfmt.FromContext(ctx)
	headers := []string{"ID", "NAME", "ADDRESS", "CITY", "COUNTRY"}
	rows := make([][]string, len(result.Data))
	for i, loc := range result.Data {
		rows[i] = []string{
			loc.ID,
			loc.Name,
			loc.Address,
			loc.City,
			loc.Country,
		}
	}

	return f.Table(headers, rows, nil)
}

func runLocationsGet(cmd *cobra.Command, args []string) error {
	locationID := args[0]

	ctx := cmd.Context()
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	location, err := client.GetLocation(ctx, threads.LocationID(locationID))
	if err != nil {
		return fmt.Errorf("failed to get location: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(locationToMap(location), jqQuery)
	}

	printLocationText(location)
	return nil
}

// locationToMap converts a Location to a map for JSON output
func locationToMap(loc *threads.Location) map[string]any {
	m := map[string]any{
		"id":   loc.ID,
		"name": loc.Name,
	}
	if loc.Address != "" {
		m["address"] = loc.Address
	}
	if loc.City != "" {
		m["city"] = loc.City
	}
	if loc.Country != "" {
		m["country"] = loc.Country
	}
	if loc.PostalCode != "" {
		m["postal_code"] = loc.PostalCode
	}
	if loc.Latitude != 0 || loc.Longitude != 0 {
		m["latitude"] = loc.Latitude
		m["longitude"] = loc.Longitude
	}
	return m
}

// printLocationText prints a Location in text format
func printLocationText(loc *threads.Location) {
	ui.Success("Location Details")
	fmt.Printf("  ID:         %s\n", loc.ID)
	fmt.Printf("  Name:       %s\n", loc.Name)
	if loc.Address != "" {
		fmt.Printf("  Address:    %s\n", loc.Address)
	}
	if loc.City != "" {
		fmt.Printf("  City:       %s\n", loc.City)
	}
	if loc.Country != "" {
		fmt.Printf("  Country:    %s\n", loc.Country)
	}
	if loc.PostalCode != "" {
		fmt.Printf("  Postal:     %s\n", loc.PostalCode)
	}
	if loc.Latitude != 0 || loc.Longitude != 0 {
		fmt.Printf("  Coords:     %.6f, %.6f\n", loc.Latitude, loc.Longitude)
	}
}
