# Exchange Data Adapter (Go)

**Component**: exchange-data-adapter-go
**Domain**: Trading, account management, order lifecycle
**Epic**: TSE-0001.5a (Exchange Account Management), TSE-0001.5b (Exchange Order Processing)
**Tech Stack**: Go, PostgreSQL, Redis
**Schema Namespace**: `exchange`

## Purpose

The Exchange Data Adapter provides data persistence services for the exchange-simulator-go component, following Clean Architecture principles. It exposes domain-driven APIs for account management, order lifecycle processing, and trade execution while abstracting database implementation details.

## Architecture Compliance

**Clean Architecture**:
- Exposes business domain concepts, not database artifacts
- Provides exchange-specific APIs tailored to trading and account management needs
- Maintains complete separation from exchange-simulator business logic
- Uses shared infrastructure with logical namespace isolation

**Domain Focus**:
- Account management with sub-account isolation
- Multi-asset balance tracking (BTC, ETH, USD, USDT)
- Order lifecycle management and matching
- Trade history and transaction audit trails
- Order books with high-frequency updates

## Data Requirements

### Account & Balance Management
- **Account Creation**: Account management system with sub-account isolation
- **Multi-Asset Balances**: Real-time balance tracking across BTC, ETH, USD, USDT
- **Account Queries**: Efficient account and balance query APIs
- **Risk Checks**: Sufficient balance validation and risk controls
- **Account Audit**: Complete audit trail for all account operations

### Order Processing
- **Order Placement**: Order placement API (initially market orders only)
- **Order Matching**: Simple order matching engine (immediate fill at market price)
- **Order Status**: Order status reporting and lifecycle management
- **Trade History**: Transaction history and execution audit trails
- **Order Books**: High-frequency order book updates and management

### Storage Patterns
- **Redis**: Active order books, real-time balances, session data, order cache
- **PostgreSQL**: Account records, trade history, audit trails, order archives

## API Design Principles

### Domain-Driven APIs
The adapter exposes trading and account concepts, not database implementation:

**Good Examples**:
```go
CreateAccount(accountSpec) -> AccountID
GetAccountPortfolio(accountID) -> Portfolio
PlaceOrder(order) -> OrderID
ProcessTrade(trade) -> TradeConfirmation
UpdateOrderBook(symbol, updates) -> OrderBookSnapshot
```

**Avoid Database Artifacts**:
```go
// Don't expose these
GetAccountTable() -> []AccountRow
UpdateBalanceRecord(id, fields) -> bool
QueryOrderHistory(sql) -> ResultSet
```

## Technology Standards

### Database Conventions
- **PostgreSQL**: snake_case for tables, columns, functions
- **Redis**: kebab-case with `exchange:` namespace prefix
- **Go**: PascalCase for public APIs, camelCase for internal functions

### Performance Requirements
- **Order Processing**: Handle high-frequency order placement and matching
- **Balance Updates**: Real-time balance updates with consistency guarantees
- **Order Book Management**: Low-latency order book updates
- **Trade Settlement**: Efficient trade processing and settlement flows

## Integration Points

### Serves
- **Primary**: exchange-simulator-go
- **Integration**: Provides trading data to risk-monitor-py and trading-system-engine-py

### Dependencies
- **Shared Infrastructure**: Single PostgreSQL + Redis instances
- **Protocol Buffers**: Via protobuf-schemas for trading and account definitions
- **Service Discovery**: Via orchestrator-docker configuration
- **Market Data**: Receives price feeds from market-data-simulator-go

## Multi-Asset Trading

### Supported Assets
- **BTC/USD**: Bitcoin trading pairs
- **ETH/USD**: Ethereum trading pairs
- **BTC/USDT**: Bitcoin to Tether pairs
- **ETH/USDT**: Ethereum to Tether pairs

### Order Types (Initial)
- **Market Orders**: Immediate execution at current market price
- **Balance Validation**: Sufficient balance checks before order acceptance
- **Risk Controls**: Basic position and balance risk validation

## Development Status

**Repository**: Active (Repository Created)
**Branch**: feature/TSE-0003.0-data-adapter-foundation
**Epic Progress**: TSE-0001.5a (Exchange Account Management) - Not Started

## Next Steps

1. **Component Configuration**: Add `.claude/` configuration for exchange-specific patterns
2. **Schema Design**: Design exchange schema in `exchange` PostgreSQL namespace
3. **API Definition**: Define account management and order processing APIs
4. **Implementation**: Implement adapter with comprehensive testing
5. **Integration**: Connect with exchange-simulator-go component

## Configuration Management

- **Shared Configuration**: project-plan/.claude/ for global architecture patterns
- **Component Configuration**: .claude/ directory for exchange-specific settings (to be created)
- **Database Schema**: `exchange` namespace with high-frequency trading optimization

---

**Epic Context**: TSE-0001 Foundation Services & Infrastructure
**Last Updated**: 2025-09-18
**Architecture**: Clean Architecture with domain-driven data persistence