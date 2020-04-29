package name

// Exchange 枚举了交易所的名称
type Exchange string

// 交易所的代号
const (
	// Special Function
	LOCAL = "LOCAL" // For local generated data

	// CryptoCurrency
	BITMEX   = "BITMEX"
	OKEX     = "OKEX"
	HUOBI    = "HUOBI"
	BITFINEX = "BITFINEX"
	BINANCE  = "BINANCE"
	BYBIT    = "BYBIT" // bybit.com
	COINBASE = "COINBASE"
	DERIBIT  = "DERIBIT"
	GATEIO   = "GATEIO"
	BITSTAMP = "BITSTAMP"

	// China
	CFFEX Exchange = "CFFEX" // China Financial Futures Exchange
	SHFE           = "SHFE"  // Shanghai Futures Exchange
	CZCE           = "CZCE"  // Zhengzhou Commodity Exchange
	DCE            = "DCE"   // Dalian Commodity Exchange
	INE            = "INE"   // Shanghai International Energy Exchange
	SSE            = "SSE"   // Shanghai Stock Exchange
	SZSE           = "SZSE"  // Shenzhen Stock Exchange
	SGE            = "SGE"   // Shanghai Gold Exchange
	WXE            = "WXE"   // Wuxi Steel Exchange
	CFETS          = "CFETS" // China Foreign Exchange Trade System
	// Global
	SMART    = "SMART"    // Smart Router for US stocks
	NYSE     = "NYSE"     // New York Stock Exchange
	NASDAQ   = "NASDAQ"   // Nasdaq Exchange
	NYMEX    = "NYMEX"    // New York Mercantile Exchange
	COMEX    = "COMEX"    // a division of theNew York Mercantile Exchange
	GLOBEX   = "GLOBEX"   // Globex of CME
	IDEALPRO = "IDEALPRO" // Forex ECN of Interactive Brokers
	CME      = "CME"      // Chicago Mercantile Exchange
	ICE      = "ICE"      // Intercontinental Exchange
	SEHK     = "SEHK"     // Stock Exchange of Hong Kong
	HKFE     = "HKFE"     // Hong Kong Futures Exchange
	HKSE     = "HKSE"     // Hong Kong Stock Exchange
	SGX      = "SGX"      // Singapore Global Exchange
	CBOT     = "CBT"      // Chicago Board of Trade
	CBOE     = "CBOE"     // Chicago Board Options Exchange
	CFE      = "CFE"      // CBOE Futures Exchange
	DME      = "DME"      // Dubai Mercantile Exchange
	EUREX    = "EUX"      // Eurex Exchange
	APEX     = "APEX"     // Asia Pacific Exchange
	LME      = "LME"      // London Metal Exchange
	BMD      = "BMD"      // Bursa Malaysia Derivatives
	TOCOM    = "TOCOM"    // Tokyo Commodity Exchange
	EUNX     = "EUNX"     // Euronext Exchange
	KRX      = "KRX"      // Korean Exchange

	OANDA = "OANDA" // oanda.com
)
