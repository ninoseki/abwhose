package matchers

import (
	"fmt"
	"strings"
)

// Matches WHOIS data to the best way to report abuse to the registrar/hosting provider.
//
// Try to keep this sorted alphabetically by ProviderName
var whoisMatchers = []matcher{
	{OnlineForm{"Cloudflare", "https://www.cloudflare.com/abuse/form"}, whoisContains("abuse@cloudflare.com")},
	{OnlineForm{"Digital Ocean", "https://www.digitalocean.com/company/contact/#abuse"}, whoisContains("descr:          Digital Ocean, Inc.")},
	{OnlineForm{"Dynadot", "https://www.dynadot.com/report_abuse.html"}, whoisContains("abuse@dynadot.com")},
	{OnlineForm{"GoDaddy", "https://supportcenter.godaddy.com/AbuseReport"}, whoisContains("abuse@godaddy.com")},
	{AbuseEmail{"Hostinger", "abuse@hostinger.com"}, whoisContains("netname:        HOSTING-SERVERS")},
	{OnlineForm{"Namecheap", "https://support.namecheap.com/index.php?/Tickets/Submit"}, whoisContains("abuse@namecheap.com")},
	{OnlineForm{"Namesilo", "https://www.namesilo.com/report_abuse.php or https://new.namesilo.com/phishing_report.php"}, whoisContains("abuse@namesilo.com")},
	{AbuseEmail{"OrangeWebsite", "abuse-dept@orangewebsite.com"}, whoisContains("abuse@orangewebsite.com")},
	{OnlineForm{"PublicDomainRegistry", "https://publicdomainregistry.com/process-for-handling-abuse/"}, whoisContains("abuse-contact@publicdomainregistry.com")},
	{OnlineForm{"Tucows", "https://tucowsdomains.com/report-abuse/"}, whoisContains("abuse@tucows.com")},
	{OnlineForm{"Tucows", "https://tucowsdomains.com/report-abuse/"}, whoisContains("domainabuse@tucows.com")},
}

func whoisContains(contents string) func(string) bool {
	return func(whois string) bool {
		return strings.Contains(whois, contents)
	}
}

func getContactsFromWHOIS(query string) ([]ProviderContact, error) {
	rawWhois, err := WHOISClient(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query whois: %s, %w", string(rawWhois), err)
	}

	var contacts []ProviderContact
	for _, m := range whoisMatchers {
		if m.matches(string(rawWhois)) {
			contacts = append(contacts, m.contact)
		}
	}

	// One of the whoisMatchers matched so return that info
	if len(contacts) > 0 {
		return contacts, nil
	}

	// Nothing matched so try and extract raw email addresses
	return getRawEmailContacts(string(rawWhois)), nil
}
