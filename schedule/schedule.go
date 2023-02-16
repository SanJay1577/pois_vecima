package schedule

import (
	"fmt"
	serviceLog "pois/config/logging"
)

type Provider string

type GetSchedule interface {
	GetProviderSchedule(string, string, string, string) (string, string)
}

type ProvidersSchedule func(string, string, string) (string, string)

var providersSchedules map[string]ProvidersSchedule

// GetProviderSchedule will get provider implementation logic by provider name to process the schedule imformation
func (p Provider) GetProviderSchedule(provider string, channel string, date string, time string) (string, string) {
	serviceLog.EsamLog.Infof("Map len : %v", len(providersSchedules))
	ps, ok := providersSchedules[provider]

	if !ok {
		serviceLog.EsamLog.Infof("No provider implementation found for %v", provider)

		ps, ok = providersSchedules["default"]

		if !ok {
			return "", "no schedule parser for provider type"
		}
	}

	returnVal, err := ps(channel, date, time)

	if err != "" {
		return "", err
	}

	return returnVal, ""
}

// registerProvider will register the provider with its implementation signature
func registerProvider(provider string, parser ProvidersSchedule) {
	if providersSchedules == nil {
		providersSchedules = make(map[string]ProvidersSchedule)
	}
	_, ok := providersSchedules[provider]
	if ok {
		fmt.Println("duplicate provider for type: ", provider)
	}

	providersSchedules[provider] = parser
}
