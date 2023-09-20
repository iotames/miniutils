package miniutils

import (
	"log"
	"strconv"
	"strings"
)

// StrToInt 字符串转整数。prune参数: 去除千分位等特殊符号。
// Examples:
// StrToInt("15,633,510", nil) -> 15633510
// StrToInt("$ 156,335,10", []string{"$", ","}) != 15633510
func StrToInt(intstr string, prune []string) int64 {
	var intRes int64 = 0
	if prune == nil {
		prune = append(prune, `,`)
	}
	for _, ostr := range prune {
		intstr = strings.TrimSpace(strings.ReplaceAll(intstr, ostr, ``))
	}
	if intstr != "" {
		intRes, _ = strconv.ParseInt(intstr, 10, 64)
	}
	return intRes
}

// GetPriceByText 价格字符串转数字(float64)
// Examples:
// GetPriceByText("$ 156,335,10.37") -> 15633510.37
// GetPriceByText("£156,335,10.37") -> 15633510.37
func GetPriceByText(priceText string) float64 {
	var price float64 = 0
	if strings.TrimSpace(priceText) == "" {
		return price
	}
	priceSplit := strings.Split(priceText, "$")
	if len(priceSplit) > 1 {
		price, _ = strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(priceSplit[1], ",", "")), 64)
	}
	if price == 0 {
		priceSplit = strings.Split(priceText, "€")
		if len(priceSplit) > 0 {
			priceStr := strings.TrimSpace(priceSplit[0])
			if priceStr == "" && len(priceSplit) > 1 {
				priceStr = strings.TrimSpace(priceSplit[1])
			}
			if priceStr != "" {
				price1 := strings.ReplaceAll(priceStr, `.`, ``)
				price2 := strings.ReplaceAll(price1, `,`, `.`)
				price, _ = strconv.ParseFloat(price2, 64)
			}
		}
	}
	if price == 0 {
		priceSplit = strings.Split(priceText, "£")
		if len(priceSplit) > 1 {
			price, _ = strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(priceSplit[1], ",", "")), 64)
		}
	}
	if price == 0 {
		priceSplit = strings.Split(priceText, `DKK`)
		log.Println("-----Check---price----DKK", priceSplit)
		if len(priceSplit) > 0 {
			price1 := strings.ReplaceAll(priceSplit[0], ".", "")
			price2 := strings.ReplaceAll(price1, `,`, `.`)
			price, _ = strconv.ParseFloat(strings.TrimSpace(price2), 64)
		}
	}
	return price
}
