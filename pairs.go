package main

type pair struct {
	first       string
	second      string
	sameMeaning bool
}

var defaultPairs = []pair{
	{"How do I reset my password?", "I forgot my password, how can I change it?", true},
	{"How do I cancel my subscription?", "I want to stop paying for my subscription", true},
	{"Which payment methods do you accept?", "What payment options are supported?", true},
	{"How do I export my data?", "Can I download all my data?", true},
	{"Where can I find my invoices?", "How do I get a copy of my bill?", true},
	{"How do I contact support?", "How can I reach customer service?", true},

	{"How do I enable 2FA?", "How do I disable 2FA?", false},
	{"How do I upgrade my plan?", "How do I downgrade my plan?", false},
	{"How do I add a team member?", "How do I remove a team member?", false},
	{"What is my current subscription plan?", "What subscription plans are available?", false},
	{"How much does it cost for 10 users?", "How much does it cost for 100 users?", false},
	{"Does it support Go 1.21?", "Does it support Go 1.22?", false},
	{"Is my data stored in the EU?", "Is my data stored in the US?", false},
	{"Why was I charged twice?", "Why wasn't I charged?", false},
}
