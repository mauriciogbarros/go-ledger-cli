# Go Transaction Ledger CLI

![Go](https://img.shields.io/badge/Go-1.22+-blue.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Contributions](https://img.shields.io/badge/contributions-welcome-orange.svg)

A command-line application written in Go that simulates a **financial transaction ledger system**, demonstrating backend engineering fundamentals such as data structures, domain modeling, and abstraction.

---

## Table of Contents
- [Problem Statement](#problem-statement)
- [Deliverable Definition](#deliverable-definition)
- [Concepts Used](#concepts-used)
- [Project Requirements](#project-requirements)
- [Core Features](#core-features)
- [Application Structure](#application-structure)
- [Development Plan](#development-plan)
- [Specifications](#specifications)
- [User Stories](#user-stories)
- [Getting Started](#getting-started)
- [Demo](#demo)
- [Portfolio Value](#portfolio-value)
- [Future Improvements](#future-improvements)
- [Contributing](#contributing)
- [License](#license)

---

## Problem Statement

Modern financial systems rely on accurate, consistent, and structured transaction processing. Even simple operations such as deposits, withdrawals, and transfers require:

- Proper data modeling
- Controlled state changes
- Clear abstractions and separation of concerns

This project addresses the problem of building a minimal but well-structured ledger system that:

- Tracks accounts and transactions
- Maintains balances correctly
- Demonstrates how core banking logic can be implemented

---

## Deliverable Definition

Build a financial transaction ledger CLI in Go demonstrating data structures, custom types, encapsulation, and interfaces.

---

## Concepts Used
### Data Structures
- Slices for storing transactions
- Maps for fast account lookup
- Structs for domain modeling

### Types & Abstraction
- Custom types (e.g., AccountID, Currency)
- Encapsulation via unexported fields
- Struct embedding for reuse
- Interfaces for decoupling behavior

### Engineering Practices
- Package-based structure
- Separation of concerns
- CLI-driven architecture

---

## Project Requirements

### Functional Requirements
- Create and manage accounts
- Record financial transactions
- Compute balances dynamically
- Track transaction history

### Non-Functional Requirements
- Clear and maintainable code structure
- Domain-driven naming
- No direct mutation of critical state
- Error handling for invalid operations

---

## Core Features
### Account Management
```bash
ledger create-account <name>
ledger list-accounts
```

### Transactions
```bash
ledger deposit <account_id> <amount>
ledger withdraw <account_id> <amount>
ledger transfer <from> <to> <amount>
```

### Ledger Operations
```bash
ledger balance <account_id>
ledger history <account_id>
```

---

## Application Structure

```
cmd/
	main.go

internal/
  account/
		account.go
  transaction/
		transaction.go
  ledger/
		ledger.go

  storage/
		memory.go

  cli/
		commands.go
```

### Structure Goals
- Separate domain logic from CLI handling
- Enable future extension (e.g., database, API)
- Keep business logic testable

---

## Development Plan
### Step 1 — Domain Modeling
- Define Account struct
- Define Transaction struct
- Define transaction types (deposit, withdraw, transfer)

### Step 2 — Data Structures
- Use slices to store transactions
- Use maps to store accounts by ID

### Step 3 — Ledger Logic
- Implement:
  - Deposit
  - Withdraw
  - Transfer
- Ensure balance consistency

### Step 4 — Abstraction
- Introduce interfaces

```go
type TransactionProcessor interface {
    Process(tx Transaction) error
}
```

```go
type Storage interface {
    SaveTransaction(tx Transaction) error
    GetTransactions(accountID string) []Transaction
}
```

### Step 5 - CLI Interface
- Parse commands
- Connect CLI input to domain logic

### Step 6 — Validation
- Prevent invalid operations:
  - Negative amounts
  - Insufficient funds
  - Introduce interfaces:

---

## Specifications
### Account
```go
type Account struct {
    ID      AccountID
    Name    string
    balance Currency
}
```

### Transaction
```go
type Transaction struct {
    ID        string
    Type      string
    Amount    Currency
    From      AccountID
    To        AccountID
    Timestamp time.Time
}
```

### Rules
- Balance must not be negative
- Every action produces a transaction
- Transfers affect two accounts

---

## User Stories
### Account Management
- As a user, I want to create an account so that I can store money
- As a user, I want to list accounts so that I can see all available accounts

### Transactions
- As a user, I want to deposit money so that I can increase my balance
- As a user, I want to withdraw money so that I can use my funds
- As a user, I want to transfer money so that I can move funds between accounts

### Ledger Visibility
- As a user, I want to see my balance so that I know how much money I have
- As a user, I want to see transaction history so that I can track activity

---

## Getting Started

```bash
git clone https://github.com/yourusername/go-ledger-cli.git
cd go-ledger-cli
go run cmd/main.go
```

---

## Demo

```bash
$ ledger create-account Alice
✅ Account created: acc_1

$ ledger create-account Bob
✅ Account created: acc_2

$ ledger deposit acc_1 1000
✅ Deposited 1000.00 to acc_1

$ ledger transfer acc_1 acc_2 250
✅ Transferred 250.00 from acc_1 to acc_2

$ ledger balance acc_1
💰 Balance: 750.00

$ ledger balance acc_2
💰 Balance: 250.00

$ ledger history acc_1
📜 Transactions:
- Deposit: 1000.00
- Transfer: -250.00 → acc_2
```

---

## Portfolio Value
This project demonstrates:
- Backend engineering fundamentals
- Financial domain modeling
- Code organization and abstraction
- Readiness for:
  - FinTech
  - Banking Systems
  - Public sector applications

---

## Future Improvements
- Persistence (file or database)
- REST API interface
- Authentication layer
- Concurrency-safe processing
- Integration with external systems

---

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

1. Fork the repository.
2. Create a new branch for your changes.
3. Make your changes and commit them.
4. Push your changes to your forked repository.
5. Open a pull request.

---

## License

This project is licensed under the [MIT License](./LICENSE.md)


# Notes
- Value conversion fallow banker's rounding
