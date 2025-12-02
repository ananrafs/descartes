# Cashback Custom Extensions - Implementation Guide

This document provides reference implementations for the custom extensions needed to support the cashback differentiation requirements.

---

## Extension 1: String Equal or All Rule

### Purpose
Implements the "ALL" wildcard logic where a field value of "ALL" matches any actual value.

### Location
`engine/rules/rule/string_equal_or_all.go`

### Implementation

```go
package rule

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "fmt"

    "github.com/your-org/descartes/engine/facts"
    "github.com/your-org/descartes/engine/rules"
    "github.com/your-org/descartes/common/parser"
)

const (
    StringEqualOrAllType = "rules.string.equal_or_all"
)

// StringEqualOrAllRule checks if a field equals a value OR if the value is "ALL" (wildcard)
type StringEqualOrAllRule struct {
    Field string `json:"field"` // The fact field to check
    Value string `json:"value"` // Expected value or "ALL" for wildcard
}

// GetType returns the rule type identifier
func (r *StringEqualOrAllRule) GetType() string {
    return StringEqualOrAllType
}

// GetHash generates a unique hash for caching
func (r *StringEqualOrAllRule) GetHash() string {
    hash := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", r.GetType(), r.Field, r.Value)))
    return hex.EncodeToString(hash[:])
}

// IsMatch evaluates the rule against facts
// Returns true if:
// 1. Value is "ALL" (wildcard matches anything), OR
// 2. Fact field value exactly matches Value
func (r *StringEqualOrAllRule) IsMatch(facts facts.FactsItf) (bool, error) {
    // AC4: "ALL" matches everything
    if r.Value == "ALL" {
        return true, nil
    }

    // Parse field from facts (supports templates)
    var actualValue string
    err := parser.DeepTemplateEvaluateFromMap(facts.GetMap(), r.Field, &actualValue)
    if err != nil {
        return false, fmt.Errorf("failed to evaluate field %s: %w", r.Field, err)
    }

    // AC4: Exact match for non-"ALL" fields
    return actualValue == r.Value, nil
}

// UnmarshalJSON implements custom JSON deserialization
func (r *StringEqualOrAllRule) UnmarshalJSON(data []byte) error {
    type Alias StringEqualOrAllRule
    aux := &struct {
        *Alias
    }{
        Alias: (*Alias)(r),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    return nil
}

// Factory implementation for registration
func (r *StringEqualOrAllRule) Create() rules.RulesItf {
    return &StringEqualOrAllRule{}
}

func (r *StringEqualOrAllRule) GetRuleType() string {
    return StringEqualOrAllType
}
```

### Unit Tests

```go
package rule

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/your-org/descartes/engine/facts"
)

func TestStringEqualOrAllRule_AllWildcard(t *testing.T) {
    rule := &StringEqualOrAllRule{
        Field: "country",
        Value: "ALL",
    }

    // Should match any value
    fact := facts.NewFactsBuilder().
        Set("country", "id").
        Build()

    match, err := rule.IsMatch(fact)
    assert.NoError(t, err)
    assert.True(t, match)

    // Should match different value
    fact2 := facts.NewFactsBuilder().
        Set("country", "sg").
        Build()

    match2, err := rule.IsMatch(fact2)
    assert.NoError(t, err)
    assert.True(t, match2)
}

func TestStringEqualOrAllRule_ExactMatch(t *testing.T) {
    rule := &StringEqualOrAllRule{
        Field: "country",
        Value: "id",
    }

    // Should match exact value
    fact := facts.NewFactsBuilder().
        Set("country", "id").
        Build()

    match, err := rule.IsMatch(fact)
    assert.NoError(t, err)
    assert.True(t, match)

    // Should NOT match different value
    fact2 := facts.NewFactsBuilder().
        Set("country", "sg").
        Build()

    match2, err := rule.IsMatch(fact2)
    assert.NoError(t, err)
    assert.False(t, match2)
}

func TestStringEqualOrAllRule_MissingField(t *testing.T) {
    rule := &StringEqualOrAllRule{
        Field: "country",
        Value: "id",
    }

    // Should return error for missing field
    fact := facts.NewFactsBuilder().
        Set("platform", "ANDROID_APP").
        Build()

    match, err := rule.IsMatch(fact)
    assert.Error(t, err)
    assert.False(t, match)
}
```

---

## Extension 2: Cashback Calculation Action

### Purpose
Implements cashback calculation with all business logic (AC5, AC6, AC7, AC10, AC11, AC13).

### Location
`engine/actions/action/cashback_calculate.go`

### Implementation

```go
package action

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/your-org/descartes/engine/actions"
    "github.com/your-org/descartes/engine/facts"
    "github.com/your-org/descartes/common/parser"
)

const (
    CashbackCalculateType = "actions.cashback.calculate"
)

// CashbackCalculationAction calculates cashback based on business rules
type CashbackCalculationAction struct {
    Name               string  `json:"name"`                  // Rule name for tracking
    CashbackRate       float64 `json:"cashback_rate"`         // Rate (0.0 to 1.0)
    MaxCashbackAmount  int     `json:"max_cashback_amount"`   // Maximum cashback cap
    ExpiryDays         int     `json:"expiry_days"`           // Days until cashback expires
    OrderAmountField   string  `json:"order_amount_field"`    // Field containing order amount
}

// CashbackResult is the output structure
type CashbackResult struct {
    MatchedRule       string    `json:"matched_rule"`
    CashbackAmount    float64   `json:"cashback_amount"`
    CashbackRate      float64   `json:"cashback_rate"`       // AC13: Store for audit
    MaxCashbackAmount int       `json:"max_cashback_amount"`
    ExpiryTime        time.Time `json:"expiry_time"`
    OrderAmount       int       `json:"order_amount"`
    Capped            bool      `json:"capped,omitempty"`
    Warning           string    `json:"warning,omitempty"`
}

// GetType returns the action type identifier
func (a *CashbackCalculationAction) GetType() string {
    return CashbackCalculateType
}

// Do executes the cashback calculation
func (a *CashbackCalculationAction) Do(facts facts.FactsItf) (interface{}, error) {
    // Extract order amount from facts
    var orderAmount int
    if err := parser.DeepTemplateEvaluateFromMap(
        facts.GetMap(),
        a.OrderAmountField,
        &orderAmount,
    ); err != nil {
        return nil, fmt.Errorf("failed to get order amount: %w", err)
    }

    // AC5: Cashback Percentage Calculation
    // cashback_amount = order_amount × cashback_rate
    cashbackAmount := float64(orderAmount) * a.CashbackRate

    result := CashbackResult{
        MatchedRule:       a.Name,
        CashbackRate:      a.CashbackRate, // AC13: Record for audit
        MaxCashbackAmount: a.MaxCashbackAmount,
        OrderAmount:       orderAmount,
        Capped:            false,
    }

    // AC10: Invalid Max Cashback Amount
    // If max_cashback_amount <= 0, granted cashback must be 0
    if a.MaxCashbackAmount <= 0 {
        result.CashbackAmount = 0
        result.Warning = "max_cashback_amount is 0, no cashback granted"
        log.Printf("[WARN] Rule %s: max_cashback_amount <= 0, setting cashback to 0", a.Name)
    } else {
        // AC6: Cashback Maximum Cap
        // If cashback_amount exceeds max_cashback_amount, cap it
        if cashbackAmount > float64(a.MaxCashbackAmount) {
            result.CashbackAmount = float64(a.MaxCashbackAmount)
            result.Capped = true
        } else {
            result.CashbackAmount = cashbackAmount
        }
    }

    // AC7 & AC11: Cashback Expiry Assignment
    currentTime := time.Now()

    if a.ExpiryDays <= 0 {
        // AC11: Invalid Expiry Days
        // Set to immediately expired and log warning
        result.ExpiryTime = currentTime
        result.Warning = "expiry_days <= 0, cashback immediately expired"
        log.Printf("[WARN] Rule %s: expiry_days <= 0, setting immediate expiry", a.Name)
    } else {
        // AC7: Normal expiry calculation
        // expiry = current_time + expiry_days
        result.ExpiryTime = currentTime.AddDate(0, 0, a.ExpiryDays)
    }

    // Optionally copy additional fields for audit trail
    if userID, ok := facts.GetMap()["user_id"]; ok {
        result.UserID = userID
    }
    if orderID, ok := facts.GetMap()["order_id"]; ok {
        result.OrderID = orderID
    }

    return result, nil
}

// UnmarshalJSON implements custom JSON deserialization
func (a *CashbackCalculationAction) UnmarshalJSON(data []byte) error {
    type Alias CashbackCalculationAction
    aux := &struct {
        *Alias
    }{
        Alias: (*Alias)(a),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    return nil
}

// Factory implementation for registration
func (a *CashbackCalculationAction) Create() actions.ActionsItf {
    return &CashbackCalculationAction{}
}

func (a *CashbackCalculationAction) GetActionType() string {
    return CashbackCalculateType
}

// Enhanced result structure with audit fields
type CashbackResultWithAudit struct {
    CashbackResult
    UserID  interface{} `json:"user_id,omitempty"`
    OrderID interface{} `json:"order_id,omitempty"`
}
```

### Unit Tests

```go
package action

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/your-org/descartes/engine/facts"
)

func TestCashbackCalculation_BasicCalculation(t *testing.T) {
    action := &CashbackCalculationAction{
        Name:              "Test Rule",
        CashbackRate:      0.10,
        MaxCashbackAmount: 50000,
        ExpiryDays:        30,
        OrderAmountField:  "order_amount",
    }

    fact := facts.NewFactsBuilder().
        Set("order_amount", 100000).
        Build()

    result, err := action.Do(fact)
    assert.NoError(t, err)

    cashback := result.(CashbackResult)
    assert.Equal(t, 10000.0, cashback.CashbackAmount) // 100000 * 0.10
    assert.Equal(t, 0.10, cashback.CashbackRate)
    assert.False(t, cashback.Capped)
}

func TestCashbackCalculation_MaxCap(t *testing.T) {
    action := &CashbackCalculationAction{
        Name:              "Test Rule",
        CashbackRate:      0.10,
        MaxCashbackAmount: 30000,
        ExpiryDays:        30,
        OrderAmountField:  "order_amount",
    }

    fact := facts.NewFactsBuilder().
        Set("order_amount", 500000).
        Build()

    result, err := action.Do(fact)
    assert.NoError(t, err)

    cashback := result.(CashbackResult)
    assert.Equal(t, 30000.0, cashback.CashbackAmount) // Capped at max
    assert.True(t, cashback.Capped)
}

func TestCashbackCalculation_InvalidMaxAmount(t *testing.T) {
    action := &CashbackCalculationAction{
        Name:              "Test Rule",
        CashbackRate:      0.10,
        MaxCashbackAmount: 0, // Invalid
        ExpiryDays:        30,
        OrderAmountField:  "order_amount",
    }

    fact := facts.NewFactsBuilder().
        Set("order_amount", 100000).
        Build()

    result, err := action.Do(fact)
    assert.NoError(t, err)

    cashback := result.(CashbackResult)
    assert.Equal(t, 0.0, cashback.CashbackAmount) // AC10: Must be 0
    assert.Contains(t, cashback.Warning, "max_cashback_amount is 0")
}

func TestCashbackCalculation_InvalidExpiryDays(t *testing.T) {
    action := &CashbackCalculationAction{
        Name:              "Test Rule",
        CashbackRate:      0.10,
        MaxCashbackAmount: 50000,
        ExpiryDays:        -5, // Invalid
        OrderAmountField:  "order_amount",
    }

    fact := facts.NewFactsBuilder().
        Set("order_amount", 100000).
        Build()

    result, err := action.Do(fact)
    assert.NoError(t, err)

    cashback := result.(CashbackResult)
    assert.WithinDuration(t, time.Now(), cashback.ExpiryTime, 1*time.Second) // AC11: Immediately expired
    assert.Contains(t, cashback.Warning, "expiry_days <= 0")
}

func TestCashbackCalculation_AuditFields(t *testing.T) {
    action := &CashbackCalculationAction{
        Name:              "Test Rule",
        CashbackRate:      0.05,
        MaxCashbackAmount: 50000,
        ExpiryDays:        30,
        OrderAmountField:  "order_amount",
    }

    fact := facts.NewFactsBuilder().
        Set("order_amount", 100000).
        Set("user_id", "user_123").
        Set("order_id", "order_456").
        Build()

    result, err := action.Do(fact)
    assert.NoError(t, err)

    cashback := result.(CashbackResult)
    assert.Equal(t, "Test Rule", cashback.MatchedRule)
    assert.Equal(t, 0.05, cashback.CashbackRate) // AC13: Rate stored for audit
}
```

---

## Extension 3: Configuration Validator

### Purpose
Validates cashback law configuration before loading (AC8, AC9).

### Location
`law/validator.go`

### Implementation

```go
package law

import (
    "encoding/json"
    "fmt"
)

// CashbackConfigValidator validates cashback law configuration
type CashbackConfigValidator struct{}

// ValidationError represents a configuration validation error
type ValidationError struct {
    RuleName string
    Field    string
    Message  string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("Validation error in rule '%s', field '%s': %s",
        e.RuleName, e.Field, e.Message)
}

// Validate checks if cashback configuration is valid
func (v *CashbackConfigValidator) Validate(lawJSON []byte) []ValidationError {
    var errors []ValidationError

    // Parse JSON
    var lawData map[string]interface{}
    if err := json.Unmarshal(lawJSON, &lawData); err != nil {
        errors = append(errors, ValidationError{
            RuleName: "N/A",
            Field:    "JSON",
            Message:  fmt.Sprintf("Invalid JSON: %v", err),
        })
        return errors
    }

    // Validate top-level structure
    evaluatorData, ok := lawData["evaluator"].(map[string]interface{})
    if !ok {
        errors = append(errors, ValidationError{
            RuleName: "N/A",
            Field:    "evaluator",
            Message:  "Missing or invalid evaluator field",
        })
        return errors
    }

    // Validate each rule in first_match evaluators
    evaluators, ok := evaluatorData["evaluators"].([]interface{})
    if !ok {
        return errors
    }

    for i, evalInterface := range evaluators {
        evaluator, ok := evalInterface.(map[string]interface{})
        if !ok {
            continue
        }

        // Extract action (contains cashback settings)
        action, ok := evaluator["action"].(map[string]interface{})
        if !ok {
            errors = append(errors, ValidationError{
                RuleName: fmt.Sprintf("Rule #%d", i),
                Field:    "action",
                Message:  "Missing action field",
            })
            continue
        }

        // Get rule name
        ruleName, _ := action["name"].(string)
        if ruleName == "" {
            ruleName = fmt.Sprintf("Rule #%d", i)
        }

        // AC8: Validate required fields
        requiredFields := []string{"cashback_rate", "max_cashback_amount", "expiry_days"}
        for _, field := range requiredFields {
            if _, exists := action[field]; !exists {
                errors = append(errors, ValidationError{
                    RuleName: ruleName,
                    Field:    field,
                    Message:  "Missing required field",
                })
            }
        }

        // AC9: Validate cashback_rate
        if rate, ok := action["cashback_rate"].(float64); ok {
            if rate < 0 || rate > 1 {
                errors = append(errors, ValidationError{
                    RuleName: ruleName,
                    Field:    "cashback_rate",
                    Message:  fmt.Sprintf("Invalid cashback_rate: %f (must be between 0 and 1)", rate),
                })
            }
        }

        // AC10: Validate max_cashback_amount (warning only, not error)
        if maxAmount, ok := action["max_cashback_amount"].(float64); ok {
            if maxAmount <= 0 {
                // This is valid per AC10, but log a warning
                fmt.Printf("[WARN] Rule %s: max_cashback_amount <= 0, will result in zero cashback\n", ruleName)
            }
        }

        // AC11: Validate expiry_days (warning only, not error)
        if expiryDays, ok := action["expiry_days"].(float64); ok {
            if expiryDays <= 0 {
                // This is valid per AC11, but log a warning
                fmt.Printf("[WARN] Rule %s: expiry_days <= 0, cashback will expire immediately\n", ruleName)
            }
        }
    }

    return errors
}

// ValidateOrPanic validates and panics if there are errors
func (v *CashbackConfigValidator) ValidateOrPanic(lawJSON []byte) {
    errors := v.Validate(lawJSON)
    if len(errors) > 0 {
        panic(fmt.Sprintf("Configuration validation failed with %d errors:\n%v", len(errors), errors))
    }
}
```

### Unit Tests

```go
package law

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestValidator_ValidConfig(t *testing.T) {
    validJSON := `{
        "slug": "cashback_evaluation",
        "evaluator": {
            "type": "evaluator.group.first_match",
            "evaluators": [
                {
                    "action": {
                        "name": "Test Rule",
                        "cashback_rate": 0.10,
                        "max_cashback_amount": 50000,
                        "expiry_days": 30
                    }
                }
            ]
        }
    }`

    validator := &CashbackConfigValidator{}
    errors := validator.Validate([]byte(validJSON))
    assert.Empty(t, errors)
}

func TestValidator_InvalidCashbackRate(t *testing.T) {
    invalidJSON := `{
        "slug": "cashback_evaluation",
        "evaluator": {
            "type": "evaluator.group.first_match",
            "evaluators": [
                {
                    "action": {
                        "name": "Test Rule",
                        "cashback_rate": 1.5,
                        "max_cashback_amount": 50000,
                        "expiry_days": 30
                    }
                }
            ]
        }
    }`

    validator := &CashbackConfigValidator{}
    errors := validator.Validate([]byte(invalidJSON))
    assert.NotEmpty(t, errors)
    assert.Contains(t, errors[0].Message, "cashback_rate")
}

func TestValidator_MissingRequiredField(t *testing.T) {
    invalidJSON := `{
        "slug": "cashback_evaluation",
        "evaluator": {
            "type": "evaluator.group.first_match",
            "evaluators": [
                {
                    "action": {
                        "name": "Test Rule",
                        "cashback_rate": 0.10
                    }
                }
            ]
        }
    }`

    validator := &CashbackConfigValidator{}
    errors := validator.Validate([]byte(invalidJSON))
    assert.NotEmpty(t, errors)
    assert.True(t, len(errors) >= 2) // Missing max_cashback_amount and expiry_days
}
```

---

## Registration and Usage

### Registering Custom Extensions

Update `core/initiator.go` to include custom extensions:

```go
package core

import (
    "github.com/your-org/descartes/engine/actions/action"
    "github.com/your-org/descartes/engine/rules/rule"
)

func InitCashbackFactory() {
    InitFactory(
        WithDefaults(),
        Factory{
            RuleCreateFunction: func() []rules.Factory {
                return []rules.Factory{
                    &rule.StringEqualOrAllRule{},
                }
            },
            ActionCreateFunction: func() []actions.Factory {
                return []actions.Factory{
                    &action.CashbackCalculationAction{},
                }
            },
        },
    )
}
```

### Loading and Validating Configuration

```go
package main

import (
    "encoding/json"
    "log"
    "os"

    "github.com/your-org/descartes/core"
    "github.com/your-org/descartes/law"
)

func main() {
    // Initialize factory with custom extensions
    core.InitCashbackFactory()

    // Load configuration
    lawJSON, err := os.ReadFile("cashback_example.json")
    if err != nil {
        log.Fatal(err)
    }

    // Validate configuration
    validator := &law.CashbackConfigValidator{}
    if errors := validator.Validate(lawJSON); len(errors) > 0 {
        log.Fatalf("Configuration validation failed: %v", errors)
    }

    // Parse and register law
    var cashbackLaw law.Law
    if err := json.Unmarshal(lawJSON, &cashbackLaw); err != nil {
        log.Fatal(err)
    }

    core.RegisterLaw(cashbackLaw)

    log.Println("Cashback law loaded and validated successfully")
}
```

---

## Testing Checklist

### Unit Tests
- [ ] String Equal or All Rule
  - [ ] "ALL" matches any value
  - [ ] Exact match works correctly
  - [ ] Missing field returns error
- [ ] Cashback Calculation Action
  - [ ] Basic calculation (AC5)
  - [ ] Max cap applied (AC6)
  - [ ] Expiry calculation (AC7)
  - [ ] Invalid max amount handled (AC10)
  - [ ] Invalid expiry days handled (AC11)
  - [ ] Audit fields included (AC13)
- [ ] Configuration Validator
  - [ ] Valid config passes
  - [ ] Invalid rate rejected (AC9)
  - [ ] Missing fields rejected (AC8)

### Integration Tests
- [ ] All 13 ACs pass end-to-end
- [ ] First-match priority works correctly
- [ ] Fallback to default rule
- [ ] Active status filtering
- [ ] Multiple rule combinations

### Performance Tests
- [ ] 1000 evaluations/second
- [ ] Cache effectiveness
- [ ] Memory usage acceptable

---

## Deployment Checklist

- [ ] Unit tests passing (100% coverage)
- [ ] Integration tests passing
- [ ] Performance benchmarks meet SLA
- [ ] Documentation updated
- [ ] Code review completed
- [ ] Configuration validation added
- [ ] Monitoring/alerting configured
- [ ] Rollback plan documented
- [ ] Stakeholder approval

---

## Support and Troubleshooting

### Common Issues

**Issue:** Rule not matching when it should
- Check field names match exactly
- Verify "ALL" logic is used correctly
- Check active status is true
- Review rule order (first-match)

**Issue:** Incorrect cashback amount
- Verify cashback_rate is between 0 and 1
- Check if capping is applied
- Ensure order_amount field is correct

**Issue:** Configuration validation fails
- Check JSON syntax
- Verify all required fields present
- Ensure cashback_rate is 0-1 range

### Debug Mode

Enable verbose logging:
```go
core.SetLogLevel(core.DEBUG)
```

This will log:
- Rule matching steps
- Cache hits/misses
- Action execution details
- Validation warnings
