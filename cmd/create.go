package cmd

import (
	"cs2itemlookup/cmd/textinput"
	"fmt"
	"net/url"
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/gocolly/colly"
)

var (
	LogoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)

func init() {
	rootCmd.AddCommand(createCmd)
}

type Options struct {
	Output *textinput.Output
}

var createCmd = &cobra.Command{
	Use:   "search",
	Short: "This is the main search command",
	Long:  ".",

	Run: func(cmd *cobra.Command, args []string) {

		options := Options{
			Output: &textinput.Output{},
		}

		tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.Output))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}

		c := colly.NewCollector()
		KnifeStarIconWithSpace := "★ "
		StatTrakIconWithSpace := "StatTrak™ "

		// Define a struct to unmarshal the JSON data
		type PriceOverview struct {
			LowestPrice string `json:"lowest_price"`
			Success bool `json:"success"`
		}

		// Construct the URL using the collected inputs
		weapon := options.Output.Weapon
		if options.Output.IsStatTrak {
			weapon = StatTrakIconWithSpace + weapon
		}
		if options.Output.IsKnife {
			weapon = KnifeStarIconWithSpace + weapon
		}
		skin := options.Output.Skin
		condition := options.Output.Condition

		marketHashName := fmt.Sprintf("%s | %s (%s)", weapon, skin, condition)

		// gotta use url.Values to encode the params properly 
		params := url.Values{}
		params.Set("appid", "730")
		params.Set("currency", "1")
		params.Set("market_hash_name", marketHashName)

		// Build the final url
		finalURL := fmt.Sprintf("https://steamcommunity.com/market/priceoverview/?%s", params.Encode())

		// Set up the OnResponse handler before visiting the URL
        c.OnResponse(func(r *colly.Response) {
            var data PriceOverview
            if err := json.Unmarshal(r.Body, &data); err != nil {
                fmt.Println("Error parsing JSON:", err)
                return
            }

            if !data.Success {
                fmt.Println("Failed to retrieve price information.")
                return
            }

            fmt.Println("Current Price:", data.LowestPrice)
        })

        // Error handling for requests
        c.OnRequest(func(r *colly.Request) {
            fmt.Println("Visiting:", r.URL.String())
        })

        c.OnError(func(_ *colly.Response, err error) {
            fmt.Println("Error occurred:", err)
        })

        // Now visit the URL
        if err := c.Visit(finalURL); err != nil {
            fmt.Println("Error visiting URL:", err)
        }
	},
}