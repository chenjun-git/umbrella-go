package lang

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	httpHeaderLanguageKey = "Accept-Language"
	defaultLanguage       = "en-US"
	metadataLanguageKey   = "language"
)

type LangQPair struct {
	Lang string
	Q    float64
}

func ParseAcceptLanguage(acceptLang string) []LangQPair {
	var results []LangQPair

	items := strings.Split(acceptLang, ",")
	fmt.Printf("%v\n", items)
	for _, langQ := range items {
		langQ = strings.Trim(langQ, " ")
		if langQ == "" {
			continue
		}
		langPair := strings.Split(langQ, ";")
		fmt.Printf("%v\n", langPair)
		if len(langPair) == 1 {
			results = append(results, LangQPair{langPair[0], 1})
		} else if len(langPair) == 2 {
			var (
				qValue float64
				err    error
			)
			qPair := strings.Split(langPair[1], "=")
			fmt.Printf("%v\n", qPair)
			if len(qPair) >= 2 {
				if qValue, err = strconv.ParseFloat(qPair[1], 64); err != nil {
					qValue = 1
				}
			} else {
				qValue = 1
			}
			results = append(results, LangQPair{langPair[0], qValue})
		}
	}
	return results
}

// 从HTTP Header中取`Accept-Language`，并将其根据Q值进行稳定排序(降序)，返回结果中位于数组前面的语言是客户端更期望的
func FromHttpHeader(header http.Header) []string {
	value := header.Get(httpHeaderLanguageKey)
	lqs := ParseAcceptLanguage(value)
	sort.SliceStable(lqs, func(i, j int) bool {
		return lqs[i].Q > lqs[j].Q
	})
	languages := make([]string, 0, len(lqs))
	for _, item := range lqs {
		languages = append(languages, item.Lang)
	}
	return languages
}

// 从Outgoing Metadata中取语言数据
func FromOutgoingContext(ctx context.Context) []string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return []string{defaultLanguage}
	}

	languages, ok := md[metadataLanguageKey]
	if !ok || len(languages) == 0 {
		return []string{defaultLanguage}
	}

	return append(languages, defaultLanguage)
}

// 从Incoming Metadata中取语言数据
func FromIncomingContext(ctx context.Context) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return []string{defaultLanguage}
	}

	languages, ok := md[metadataLanguageKey]
	if !ok || len(languages) == 0 {
		return []string{defaultLanguage}
	}

	return append(languages, defaultLanguage)
}

// 对context的outgoing metadata填充languages
func ContextSetLanguages(ctx context.Context, languages []string) context.Context {
	keyValues := make([]string, 0, len(languages)*2)
	for _, v := range languages {
		keyValues = append(keyValues, metadataLanguageKey, v)
	}
	metadataNew := metadata.Pairs(keyValues...)
	metadataOld, _ := metadata.FromOutgoingContext(ctx)
	return metadata.NewOutgoingContext(ctx, metadata.Join(metadataOld, metadataNew))
}
