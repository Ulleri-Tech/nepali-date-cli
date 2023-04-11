/*
Copyright © 2023 Binij Shrestha <hello@binijshrestha.com.np>
*/
package cmd

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Ulleri-Tech_nepali-date-cli",
	Short: "Convert AD Date to BS Date",
	Long: `A CLI to Convert (Gregorian Calender) AD Date to (Bikram Sambhat) BS Date.
	For example:
		convertdate -t 
		convertdate -A 2023-04-11
		convertdate -B 2079-12-28
	`,
	Run: func(cmd *cobra.Command, args []string) {
		flagToday, err := cmd.Flags().GetBool("today")
		if err != nil {
			fmt.Println("Error getting --convert flag")
			return
		}
		flagAD, err := cmd.Flags().GetBool("ad")
		if err != nil {
			fmt.Println("Error getting --convert flag")
			return
		}
		flagBS, err := cmd.Flags().GetBool("bs")
		if err != nil {
			fmt.Println("Error getting --convert flag")
			return
		}

		if flagBS && flagAD {
			fmt.Println("Error: Both --bs or -B and --ad or -A flags cannot be used together.")
			return
		}

		if flagToday {
			todayAD := time.Now().Format("2006-01-02")
			t, err := time.Parse("2006-01-02", todayAD)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				return
			}
			convertdate, err := ADTOBS(t.Year(), int(t.Month()), t.Day())
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(convertdate, "BS")
		}

		if flagAD {
			dateStr := args[0]
			if !validateDateArg(dateStr) {
				fmt.Println("Error: Invalid date format. Please provide a date in the YYYY-MM-DD format.")
				return
			}
			t, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				return
			}
			convertdate, err := ADTOBS(t.Year(), int(t.Month()), t.Day())
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(convertdate, "BS")
		}
		if flagBS {
			dateStr := args[0]
			if !validateDateArg(dateStr) {
				fmt.Println("Error: Invalid date format. Please provide a date in the YYYY-MM-DD format.")
				return
			}

			year, month, day, err := splitDate(dateStr)
			if err != nil {
				fmt.Println("Error splitting date:", err)
				return
			}

			convertdate, err := BSTOAD(int(year), month, day)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(convertdate, "AD")
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("ad", "A", false, "Convert date from AD to BS")
	rootCmd.Flags().BoolP("bs", "B", false, "Convert date from BS to AD")
	rootCmd.Flags().BoolP("today", "t", false, "Today's date of Bikram Sambhat (BS) Calender")
}

func ADTOBS(year int, month int, day int) (string, error) {
	startDate := time.Date(MIN_AD_YEAR, time.Month(MIN_AD_MONTH+1), MIN_AD_DAY, 0, 0, 0, 0, time.UTC)
	if year == 0 {
		year = time.Now().Year()
	}
	if month == 0 {
		month = int(time.Now().Month())
	}
	if day == 0 {
		day = time.Now().Day()
	}

	dayDiff := int(math.Ceil(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Sub(startDate).Hours() / 24))

	if dayDiff < 0 || dayDiff > MAX_DAY_DIFF {
		return "", fmt.Errorf("invalid date out of range")
	}
	if month < 1 || month > 12 || day < 0 || day > 31 {
		return "", fmt.Errorf("invalid month or day")
	}
	year_bs := strconv.Itoa(MIN_BS_YEAR)
	totalCount := 0
	countUptoPreviousYear := 0

	// Find Year only first
	for currYear := MIN_BS_YEAR; totalCount < dayDiff && currYear < MAX_BS_YEAR; currYear = currYear + 1 {
		countUptoPreviousYear = totalCount
		totalCount += numberOfDaysEachYearBS[strconv.Itoa(currYear)]
		year_bs = strconv.Itoa(currYear)
	}

	nthDayofYear := dayDiff - countUptoPreviousYear + 1
	//find month
	total_count_month_wise := 0
	countUptoPreviousMonth := 0
	month_bs := 0

	for _, days := range numberOfDaysEachMonthBS[year_bs] {
		if nthDayofYear <= total_count_month_wise {
			break
		}
		countUptoPreviousMonth = total_count_month_wise
		total_count_month_wise += days
		if countUptoPreviousMonth > 0 {
			month_bs++
		}
	}

	day_bs := nthDayofYear - countUptoPreviousMonth

	return fmt.Sprintf("%s-%02d-%02d", year_bs, month_bs+1, day_bs), nil
}

func BSTOAD(year_bs int, month_bs int, day_bs int) (string, error) {

	if year_bs < MIN_BS_YEAR || year_bs > MAX_BS_YEAR {
		return "", fmt.Errorf("invalid date out of range")
	}
	if month_bs < 1 || month_bs > 12 || day_bs < 1 || day_bs > 32 {
		return "", fmt.Errorf("invalid month or day")
	}
	dayCount := 0

	// Day counting upto Prev Year
	for currYear := MIN_BS_YEAR; currYear < year_bs && currYear < MAX_BS_YEAR; currYear += 1 {
		if year_bs == MIN_BS_YEAR {
			break
		}
		dayCount += numberOfDaysEachYearBS[strconv.Itoa(currYear)]
	}

	// Day counting of This Year
	for i := 0; i < month_bs-1; i++ {
		dayCount += numberOfDaysEachMonthBS[strconv.Itoa(year_bs)][i]
	}
	dayCount += day_bs - 1

	startTime := time.Date(MIN_AD_YEAR, time.Month(MIN_AD_MONTH+1), MIN_AD_DAY, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)

	finalTimeMillis := startTime + (int64(dayCount * MILLS_IN_DAY))

	return time.Unix(finalTimeMillis/1000, 0).Format("2006-01-02"), nil
}

func validateDateArg(dateStr string) bool {
	pattern := `^\d{4}-\d{2}-\d{2}$`
	match, err := regexp.MatchString(pattern, dateStr)
	if err != nil {
		return false
	}
	return match
}

func splitDate(dateStr string) (year, month, day int, err error) {
	dateParts := strings.Split(dateStr, "-")

	// Extract the year, month, and day components and convert them to integers
	year, err = strconv.Atoi(dateParts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid year: %s", dateParts[0])
	}

	month, err = strconv.Atoi(dateParts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid month: %s", dateParts[1])
	}

	day, err = strconv.Atoi(dateParts[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid day: %s", dateParts[2])
	}

	return year, month, day, nil
}
