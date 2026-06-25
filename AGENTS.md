# Descartes Repository Overview

## Project Summary

Descartes is a tool to evaluate facts against rules-action rulesets (called "laws"). It provides a framework for defining logical rules and evaluating input data (facts) to produce results based on those rules.

## Key Features

- **Laws**: Define rule sets with slugs as identifiers and evaluators to process facts
- **Evaluators**: Translate facts to results (supports multiple evaluator types)
- **Rules**: Validate conditions before executing actions
- **Actions**: Modify data based on evaluated conditions
- **Extensible Architecture**: Add custom evaluators, rules, and actions
- **Cache Support**: Optimize repeated evaluations

## Directory Structure

```
descartes/
├── bench/             # Benchmark tests
├── cache/             # Cache implementations
├── common/            # Common utilities (parsers, converters, type checkers)
├── core/              # Core functionality (initiator, evaluation)
├── dump/              # Example law and fact JSON files with outputs
├── engine/            # Main engine
│   ├── actions/       # Action implementations
│   ├── evaluators/    # Evaluator implementations
│   ├── facts/         # Fact handling
│   └── rules/         # Rule implementations
├── errors/            # Error handling
├── law/               # Law and fact types
├── go.mod/go.sum      # Go modules
├── main.go            # CLI main file
├── Makefile           # Build file
└── README.md          # Project documentation
```

## Core Concepts

### Law
A Law is a rule set with a slug (identifier) and an evaluator.

### Evaluator
An evaluator translates facts into results. Default evaluators:
- `evaluator`: Basic evaluator with rule and action
- `evaluator.group.first_match`: Evaluates first matching rule
- `evaluator.group.multi_match`: Evaluates all matching rules
- `evaluator.group.multi_match_ordered`: Evaluates ordered matching rules
- `evaluator.group.multi_match_ordered_cycle`: Cycles through matching rules
- `evaluator.iter`: Iterates over data

### Rules
Rules validate conditions. Default rule types include:
- Integer comparisons (equal, greater, lesser, between)
- String comparisons (equal, equal_fold, wildcard)
- Boolean rules
- Array rules (contains, in, struct_contains)
- Logical operators (and, or, not)
- Time rules (before, after, between)
- Default and oneof rules

### Actions
Actions modify data. Default actions:
- Integer operations (sum, subtract, multiply, divide, mod)
- Float operations (sum, subtract, multiply, divide)
- Array operations (extend_each)
- Map operations (append)
- Action groups for multiple actions

## Usage

### Initialization
```go
core.InitFactory(core.WithDefaults())
```

### Create and Register a Law
```go
l, err := law.CreateLaw(jsonLaw)
if err != nil {
    panic(err)
}
err = core.Register(l)
```

### Create a Fact
```go
fact, err := law.CreateFact(jsonFact)
```

### Evaluate Fact
```go
res, err := core.Eval(fact)
```

## CLI Usage

```bash
go run main.go -folder ./dump/law_2 -fact fact -law law -out output
```

## Technology Stack

- **Language**: Go 1.22
- **JSON Library**: github.com/json-iterator/go
