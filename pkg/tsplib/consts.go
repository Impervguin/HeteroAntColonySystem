package tsplib

const (
	// Problem types
	TypeTSP  = "TSP"
	TypeATSP = "ATSP"

	// Edge weight types
	WeightTypeEUC2D    = "EUC_2D"
	WeightTypeEUC3D    = "EUC_3D"
	WeightTypeMAX2D    = "MAX_2D"
	WeightTypeMAX3D    = "MAX_3D"
	WeightTypeMAN2D    = "MAN_2D"
	WeightTypeMAN3D    = "MAN_3D"
	WeightTypeCEIL2D   = "CEIL_2D"
	WeightTypeGEO      = "GEO"
	WeightTypeATT      = "ATT"
	WeightTypeEXPLICIT = "EXPLICIT"

	// Edge weight formats
	WeightFormatFUNCTION       = "FUNCTION"
	WeightFormatFULL_MATRIX    = "FULL_MATRIX"
	WeightFormatUPPER_ROW      = "UPPER_ROW"
	WeightFormatLOWER_ROW      = "LOWER_ROW"
	WeightFormatUPPER_DIAG_ROW = "UPPER_DIAG_ROW"
	WeightFormatLOWER_DIAG_ROW = "LOWER_DIAG_ROW"

	// Data sections
	SectionNodeCoord   = "NODE_COORD_SECTION"
	SectionEdgeWeight  = "EDGE_WEIGHT_SECTION"
	SectionDisplayData = "DISPLAY_DATA_SECTION"
	SectionEOF         = "EOF"
)
