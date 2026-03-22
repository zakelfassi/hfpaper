package main

// getArgs safely extracts the arguments map from an MCP CallToolRequest
func getArgs(arguments interface{}) map[string]interface{} {
	if m, ok := arguments.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
}

func getStringArg(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}

func getIntArg(args map[string]interface{}, key string, defaultVal int) int {
	if v, ok := args[key].(float64); ok {
		return int(v)
	}
	return defaultVal
}

func getBoolArg(args map[string]interface{}, key string) bool {
	if v, ok := args[key].(bool); ok {
		return v
	}
	return false
}
