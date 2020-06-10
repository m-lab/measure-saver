package model

// Measurement is the data format exchanged between the uploader and this
// service.
type Measurement struct {
	ID int64 `pg:",pk"`

	// BrowserID is a unique identifier for this Measurement's uploader.
	BrowserID string `sql:",notnull"`

	// DeviceType is a free-form string identifying the device type.
	DeviceType string

	// Notes is a free-form string containing browser-specific notes.
	Notes string

	Download float64 `pg:",use_zero"`
	Upload   float64 `pg:",use_zero"`
	Latency  int     `pg:",use_zero"`

	Results Results
}

// Results represents the NDT variables sent by the server at the end of
// a measurement.
type Results struct {
	AckPktsIn            string
	ClientToServerSpeed  float64
	CongestionSignals    string
	CountRTT             string
	CurMSS               string
	CurRTO               string
	DataBytesOut         string
	DupAcksIn            string
	Jitter               float64
	MaxCwnd              string
	MaxRTT               string
	MaxRwinRcvd          string
	MinRTT               string
	S2CClientIP          string `json:"NDTResult.S2C.ClientIP"`
	S2CClientPort        string `json:"NDTResult.S2C.ClientPort"`
	S2CEndTime           string `json:"NDTResult.S2C.EndTime"`
	S2CError             string `json:"NDTResult.S2C.Error"`
	S2CMinRTT            string `json:"NDTResult.S2C.MinRTT"`
	S2CServerIP          string `json:"NDTResult.S2C.ServerIP"`
	S2CServerPort        string `json:"NDTResult.S2C.ServerPort"`
	S2CStartTime         string `json:"NDTResult.S2C.StartTime"`
	S2CUUID              string `json:"NDTResult.S2C.UUID"`
	PktsOut              string
	PktsRetrans          string
	RcvWinScale          string
	ServerToClientSpeed  float64
	SndLimTimeCwnd       string
	SndLimTimeRwin       string
	SndLimTimeSender     string
	SndWinScale          string
	Sndbuf               string
	SumRTT               string
	TCPInfoATO           string `json:"TCPInfo.ATO"`
	TCPInfoAdvMSS        string `json:"TCPInfo.AdvMSS"`
	TCPInfoAppLimited    string `json:"TCPInfo.AppLimited"`
	TCPInfoBackoff       string `json:"TCPInfo.Backoff"`
	TCPInfoBusyTime      string `json:"TCPInfo.BusyTime"`
	TCPInfoBytesAcked    string `json:"TCPInfo.BytesAcked"`
	TCPInfoBytesReceived string `json:"TCPInfo.BytesReceived"`
	TCPInfoBytesRetrans  string `json:"TCPInfo.BytesRetrans"`
	TCPInfoBytesSent     string `json:"TCPInfo.BytesSent"`
	TCPInfoCAState       string `json:"TCPInfo.CAState"`
	TCPInfoDSackDups     string `json:"TCPInfo.DSackDups"`
	TCPInfoDataSegsIn    string `json:"TCPInfo.DataSegsIn"`
	TCPInfoDataSegsOut   string `json:"TCPInfo.DataSegsOut"`
	TCPInfoDelivered     string `json:"TCPInfo.Delivered"`
	TCPInfoDeliveredCE   string `json:"TCPInfo.DeliveredCE"`
	TCPInfoDeliveryRate  string `json:"TCPInfo.DeliveryRate"`
	TCPInfoFackets       string `json:"TCPInfo.Fackets"`
	TCPInfoLastAckRecv   string `json:"TCPInfo.LastAckRecv"`
	TCPInfoLastAckSent   string `json:"TCPInfo.LastAckSent"`
	TCPInfoLastDataRecv  string `json:"TCPInfo.LastDataRecv"`
	TCPInfoLastDataSent  string `json:"TCPInfo.LastDataSent"`
	TCPInfoLost          string `json:"TCPInfo.Lost"`
	TCPInfoMaxPacingRate string `json:"TCPInfo.MaxPacingRate"`
	TCPInfoMinRTT        string `json:"TCPInfo.MinRTT"`
	TCPInfoNotsentBytes  string `json:"TCPInfo.NotsentBytes"`
	TCPInfoOptions       string `json:"TCPInfo.Options"`
	TCPInfoPMTU          string `json:"TCPInfo.PMTU"`
	TCPInfoPacingRate    string `json:"TCPInfo.PacingRate"`
	TCPInfoProbes        string `json:"TCPInfo.Probes"`
	TCPInfoRTO           string `json:"TCPInfo.RTO"`
	TCPInfoRTT           string `json:"TCPInfo.RTT"`
	TCPInfoRTTVar        string `json:"TCPInfo.RTTVar"`
	TCPInfoRWndLimited   string `json:"TCPInfo.RWndLimited"`
	TCPInfoRcvMSS        string `json:"TCPInfo.RcvMSS"`
	TCPInfoRcvRTT        string `json:"TCPInfo.RcvRTT"`
	TCPInfoRcvSpace      string `json:"TCPInfo.RcvSpace"`
	TCPInfoRcvSsThresh   string `json:"TCPInfo.RcvSsThresh"`
	TCPInfoReordSeen     string `json:"TCPInfo.ReordSeen"`
	TCPInfoReordering    string `json:"TCPInfo.Reordering"`
	TCPInfoRetrans       string `json:"TCPInfo.Retrans"`
	TCPInfoRetransmits   string `json:"TCPInfo.Retransmits"`
	TCPInfoSacked        string `json:"TCPInfo.Sacked"`
	TCPInfoSegsIn        string `json:"TCPInfo.SegsIn"`
	TCPInfoSegsOut       string `json:"TCPInfo.SegsOut"`
	TCPInfoSndBufLimited string `json:"TCPInfo.SndBufLimited"`
	TCPInfoSndCwnd       string `json:"TCPInfo.SndCwnd"`
	TCPInfoSndMSS        string `json:"TCPInfo.SndMSS"`
	TCPInfoSndSsThresh   string `json:"TCPInfo.SndSsThresh"`
	TCPInfoState         string `json:"TCPInfo.State"`
	TCPInfoTotalRetrans  string `json:"TCPInfo.TotalRetrans"`
	TCPInfoUnacked       string `json:"TCPInfo.Unacked"`
	TCPInfoWScale        string `json:"TCPInfo.WScale"`
	Timeouts             string
	c2sRate              float64
	s2cRate              float64
}
