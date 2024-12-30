package gecharts

// https://echarts.apache.org/examples/zh/editor.html?c=pie-simple
func Pie(title string, data map[string]any) map[string]any {
	sv := []map[string]any{}
	for k, v := range data {
		sv = append(sv, map[string]any{
			"value": v,
			"name":  k,
		})
	}
	return map[string]any{
		"title": map[string]any{
			"text": title,
			"left": "center",
		},
		"tooltip": map[string]any{"trigger": "item"},
		"legend": map[string]any{
			"orient": "vertical",
			"left":   "left",
		},
		"series": []map[string]any{
			{
				"name":   title,
				"type":   "pie",
				"radius": "50%",
				"emphasis": map[string]any{
					"itemStyle": map[string]any{
						"shadowBlur":    10,
						"shadowOffsetX": 0,
						"shadowColor":   "rgba(0, 0, 0, 0.5)",
					},
				},
				"data": sv,
			},
		},
	}
}
