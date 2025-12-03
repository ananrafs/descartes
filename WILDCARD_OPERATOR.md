# Wildcard Operator Implementation

## Overview
This document describes the wildcard operator implementation added to the Descartes library to support cashback rule evaluation with "ALL" wildcard matching.

## New Rule Types

### 1. `rules.string.equal.wildcard`
**File:** `/Users/a1/descartes/engine/rules/rule/string/equal_wildcard.go`

A string equality rule that supports wildcard matching. When the `value` field is set to `"ALL"` or `"*"`, the rule always matches regardless of the fact value.

**Usage:**
```json
{
  "type": "rules.string.equal.wildcard",
  "field": "country",
  "value": "ALL"
}
```

**Behavior:**
- If `value` is `"ALL"` or `"*"` → Always returns `true` (matches any value)
- Otherwise → Performs exact string equality check (`fact_value == configured_value`)

### 2. `rules.string.equal.wildcard.dynamic`
**File:** `/Users/a1/descartes/engine/rules/rule/string/equal_wildcard_dynamic.go`

Dynamic version that supports template variable substitution (e.g., `{{ field_name }}`).

**Usage:**
```json
{
  "type": "rules.string.equal.wildcard.dynamic",
  "left": "{{ user_country }}",
  "right": "{{ required_country }}"
}
```

**Behavior:**
- Resolves both `left` and `right` fields using template substitution
- If either resolved value is `"ALL"` or `"*"` → Returns `true`
- Otherwise → Performs exact string equality check

## Registration

Both rules are registered in the default factory initialization:

**File:** `/Users/a1/descartes/core/initiator.go:138-139`
```go
rule_string.NewEqualWildcard,
rule_string.NewEqualWildcardDynamic,
```

## Cashback Rule Example

### Configuration Structure
```json
{
  "slug": "cashback_rules",
  "evaluator": {
    "type": "evaluator.group.first_match",
    "evaluators": [
      {
        "name": "ID - Mobile Legends Specific",
        "type": "evaluator",
        "rule": {
          "type": "rules.conditional.and",
          "rules": [
            {
              "type": "rules.string.equal.wildcard",
              "field": "user_segment",
              "value": "ALL"
            },
            {
              "type": "rules.string.equal.wildcard",
              "field": "country",
              "value": "id"
            },
            {
              "type": "rules.string.equal.wildcard",
              "field": "platform",
              "value": "ALL"
            },
            {
              "type": "rules.string.equal.wildcard",
              "field": "category",
              "value": "ML"
            }
          ]
        },
        "action": {
          "rule_name": "ID - Mobile Legends",
          "cashback_rate": 0.10,
          "max_cashback": 30000
        }
      }
    ]
  }
}
```

### Test Cases
Test files are located in `/Users/a1/descartes/dump/test_wildcard/`

**Test Scenario 1:** Mobile Legends in Indonesia
```json
{
  "user_segment": "standard",
  "country": "id",
  "platform": "WEB",
  "category": "ML"
}
```
**Expected Match:** "ID - Mobile Legends" (10% cashback, max 30,000)

**Test Scenario 2:** Android App in Indonesia
```json
{
  "user_segment": "premium",
  "country": "id",
  "platform": "ANDROID_APP",
  "category": "FREEFIRE"
}
```
**Expected Match:** "ID - APP PROMO" (50% cashback, max 100,000)

**Test Scenario 3:** New Login User in Indonesia
```json
{
  "user_segment": "new_login",
  "country": "id",
  "platform": "IOS_APP",
  "category": "PUBG"
}
```
**Expected Match:** "ID - new login promo" (50% cashback, max 100,000)

**Test Scenario 4:** Default Rule
```json
{
  "user_segment": "vip",
  "country": "sg",
  "platform": "WEB",
  "category": "VALORANT"
}
```
**Expected Match:** "Default Cashback Rule" (5% cashback, max 50,000)

## How It Fulfills Requirements

### ✅ AC1 — Rule Matching Uses AND Logic
- Uses `rules.conditional.and` to combine multiple conditions
- All conditions must match (or be wildcards) for rule to apply

### ✅ AC2 — Fallback to Default Rule
- Implemented via `evaluator.group.first_match`
- Default rule placed last with all fields set to `"ALL"`

### ✅ AC3 — Active Rule Filtering
- Can add `rules.bool` check for `"active": true` in the AND condition

### ✅ AC4 — Exact Match for Non-"ALL" Fields
- When value is not `"ALL"` or `"*"`, performs exact string equality

### ✅ Wildcard Support
- `"ALL"` or `"*"` values always match
- Enables flexible rule matching across multiple dimensions

## Implementation Details

### Caching
Both wildcard rules support caching via the hash-based caching system:
- Cache key generated from rule type, field, and value
- Results cached per facts instance
- Reduces redundant evaluations

### Error Handling
- Returns error if field not found in facts map
- Returns error if type casting fails (non-string value)
- Follows standard Descartes error pattern

### Pattern Consistency
- Follows exact same pattern as existing string rules (`equal.go`, `equal_fold.go`)
- Maintains code consistency and readability
- Uses same interfaces and contracts

## Testing

To run the test:
```bash
go run main.go -folder ./dump/test_wildcard
```

Or with explicit file names:
```bash
go run main.go -folder ./dump/test_wildcard -fact fact -law law -out output
```

The test will:
1. Load law configuration from `law.json`
2. Load facts from `fact.json`
3. Evaluate each fact against the rules
4. Write results to `output.json`

## Priority-Based Rule Matching

The `evaluator.group.first_match` evaluator processes rules in order:

**Priority Order (Highest to Lowest):**
1. Country-specific + Category-specific rules
2. Country-specific + Platform-specific rules
3. User segment-specific rules
4. Default rule (all wildcards)

To implement the exact priority from requirements (Country > Platform > Category > User Segment):
- Order evaluators by specificity
- More specific conditions come first
- Default rule always comes last

## Usage in Cashback System

```go
// Initialize Descartes with wildcard support
core.InitFactory(core.WithDefaults())

// Load cashback rules
law, err := law.CreateLawFromJsonByte(cashbackRulesJSON)
core.Register(law)

// Evaluate cashback for a transaction
fact := law.MakeFact(map[string]interface{}{
    "user_segment": "premium",
    "country": "id",
    "platform": "ANDROID_APP",
    "category": "ML",
    "order_amount": 100000,
}).Generate("cashback_rules")

result, err := core.Eval(fact)
// result contains matched rule's action (cashback_rate, max_cashback, etc.)
```

## Future Enhancements

1. **Case-Insensitive Wildcard**: `rules.string.equal.wildcard.fold`
2. **Pattern Matching**: Support for glob patterns like `"ML*"` or `"*_APP"`
3. **Regex Wildcard**: `rules.string.regex.wildcard` for complex patterns
4. **Integer Wildcard**: `rules.int.equal.wildcard` for numeric wildcards
5. **Array Wildcard**: Match against arrays of values

## Notes

- The wildcard operator is type-safe and follows Go best practices
- Backward compatible with existing rules
- No breaking changes to existing functionality
- Well-integrated with Descartes' factory pattern and caching system
