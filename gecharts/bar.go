package gecharts

// https://echarts.apache.org/examples/zh/editor.html?c=bar-tick-align
func BarBase(title string, data map[string][]any, axixs []string, isSaveImg bool) map[string]any {
	legend := []string{}
	series := []map[string]any{}
	for k, item := range data {
		legend = append(legend, k)
		series = append(series, map[string]any{
			"name":  k,
			"type":  "bar",
			"stack": "Total",
			"data":  item,
		})
	}
	ret := map[string]any{
		"title": map[string]string{"text": title},
		"tooltip": map[string]any{
			"trigger": "axis",
			"axisPointer": map[string]any{
				"type":  "cross",
				"label": map[string]string{"backgroundColor": "#6a7985"},
			},
		},
		"legend": map[string]any{"data": legend},
		"grid":   map[string]any{"left": "3%", "right": "4%", "bottom": "3%", "containLabel": true},
		"xAxis": map[string]any{
			"type":        "category",
			"boundaryGap": false,
			"data":        axixs,
		},
		"yAxis":  map[string]any{"type": "value"},
		"series": series,
	}
	if isSaveImg {
		ret["toolbox"] = map[string]any{"feature": map[string]any{"saveAsImage": map[string]any{}}}
	}
	return ret
}

// https://echarts.apache.org/examples/zh/editor.html?c=bar-y-category
func BarHBase(title string, data map[string][]any, axixs []string, isSaveImg bool) map[string]any {
	legend := []string{}
	series := []map[string]any{}
	for k, item := range data {
		legend = append(legend, k)
		series = append(series, map[string]any{
			"name":  k,
			"type":  "bar",
			"stack": "Total",
			"data":  item,
		})
	}
	ret := map[string]any{
		"title": map[string]string{"text": title},
		"tooltip": map[string]any{
			"trigger": "axis",
			"axisPointer": map[string]any{
				"type":  "cross",
				"label": map[string]string{"backgroundColor": "#6a7985"},
			},
		},
		"legend": map[string]any{"data": legend},
		"grid":   map[string]any{"left": "3%", "right": "4%", "bottom": "3%", "containLabel": true},
		"yAxis": map[string]any{
			"type":        "category",
			"boundaryGap": false,
			"data":        axixs,
		},
		"xAxis":  map[string]any{"type": "value"},
		"series": series,
	}
	if isSaveImg {
		ret["toolbox"] = map[string]any{"feature": map[string]any{"saveAsImage": map[string]any{}}}
	}
	return ret
}

func BarStached(title string, data map[string][]any, axixs []string, isSaveImg bool) map[string]any {

	series := []map[string]any{}
	for k, item := range data {
		series = append(series, map[string]any{
			"name":     k,
			"type":     "bar",
			"stack":    "total",
			"data":     item,
			"label":    map[string]any{"show": true},
			"emphasis": map[string]any{"focus": "series"},
		})
	}
	ret := map[string]any{
		"title": map[string]string{"text": title},
		"tooltip": map[string]any{
			"trigger": "axis",
			"axisPointer": map[string]any{
				"type":  "cross",
				"label": map[string]string{"backgroundColor": "#6a7985"},
			},
		},
		"legend": map[string]any{},
		"grid":   map[string]any{"left": "3%", "right": "4%", "bottom": "3%", "containLabel": true},
		"xAxis": map[string]any{
			"type":        "category",
			"boundaryGap": false,
			"data":        axixs,
		},
		"yAxis":  map[string]any{"type": "value"},
		"series": series,
	}
	if isSaveImg {
		ret["toolbox"] = map[string]any{"feature": map[string]any{"saveAsImage": map[string]any{}}}
	}
	return ret
}
func BarHStached(title string, data map[string][]any, axixs []string, isSaveImg bool) map[string]any {

	series := []map[string]any{}
	for k, item := range data {
		series = append(series, map[string]any{
			"name":     k,
			"type":     "bar",
			"stack":    "total",
			"data":     item,
			"label":    map[string]any{"show": true},
			"emphasis": map[string]any{"focus": "series"},
		})
	}
	ret := map[string]any{
		"title": map[string]string{"text": title},
		"tooltip": map[string]any{
			"trigger": "axis",
			"axisPointer": map[string]any{
				"type":  "cross",
				"label": map[string]string{"backgroundColor": "#6a7985"},
			},
		},
		"legend": map[string]any{},
		"grid":   map[string]any{"left": "3%", "right": "4%", "bottom": "3%", "containLabel": true},
		"yAxis": map[string]any{
			"type":        "category",
			"boundaryGap": false,
			"data":        axixs,
		},
		"xAxis":  map[string]any{"type": "value"},
		"series": series,
	}
	if isSaveImg {
		ret["toolbox"] = map[string]any{"feature": map[string]any{"saveAsImage": map[string]any{}}}
	}
	return ret
}
