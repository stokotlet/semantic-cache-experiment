# semantic-cache-experiment

A semantic cache for LLM answers embeds each question and returns a cached
answer if a new question is similar enough to an old one. I wanted to check
if there is a similarity threshold that catches paraphrases but not questions
that look alike and mean something else.

The pairs are in [`pairs.go`](pairs.go):

- `SAME` - paraphrases, the cached answer is correct
- `TRAP` - look-alike questions, the cached answer is wrong

## Results

`text-embedding-3-small`, cosine similarity:

```text
SIMILARITY  TYPE  FIRST                                  SECOND
0.935       TRAP  Does it support Go 1.21?               Does it support Go 1.22?
0.852       TRAP  How do I enable 2FA?                   How do I disable 2FA?
0.838       TRAP  Is my data stored in the EU?           Is my data stored in the US?
0.834       TRAP  How much does it cost for 10 users?    How much does it cost for 100 users?
0.812       TRAP  How do I upgrade my plan?              How do I downgrade my plan?
0.780       TRAP  What is my current subscription plan?  What subscription plans are available?
0.767       SAME  How do I reset my password?            I forgot my password, how can I change it?
0.738       SAME  Which payment methods do you accept?   What payment options are supported?
0.714       SAME  How do I cancel my subscription?       I want to stop paying for my subscription
0.699       SAME  How do I contact support?              How can I reach customer service?
0.696       TRAP  Why was I charged twice?               Why wasn't I charged?
0.639       TRAP  How do I add a team member?            How do I remove a team member?
0.612       SAME  How do I export my data?               Can I download all my data?
0.540       SAME  Where can I find my invoices?          How do I get a copy of my bill?

THRESHOLD  SAME CAUGHT (useful hits)  TRAPS CAUGHT (wrong answers)
0.70       3/6                        6/8
0.75       1/6                        6/8
0.80       0/6                        5/8
0.85       0/6                        2/8
0.90       0/6                        1/8
0.95       0/6                        0/8

Median embedding latency, 1 question, 5 samples: 413ms
```

The 6 most similar pairs are all traps, so there is no threshold that works.
Paraphrases use different words and score low, while look-alikes differ by one
word and score high.

Latency was measured from my laptop, so it depends a lot on the network.

14 pairs is not a benchmark. Bigger models or a reranker would probably do
better.

## Run

```bash
export OPENAI_API_KEY=sk-...
go run .
go run . -model text-embedding-3-large
```

Tests don't need a key:

```bash
go test ./...
```
