package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
)

// JSONData represents the structure of your JSON data
type JSONData struct {
	Result struct {
		Data struct {
			Response struct {
				URL        string `json:"url"`
				Title      string `json:"title"`
				Redirects  []struct {
					Headers     struct {
						Location   string `json:"location"`
						Connection string `json:"connection"`
					} `json:"headers"`
					Status      struct {
						Code int `json:"code"`
					} `json:"status"`
					RedirectURI string `json:"redirect_uri"`
				} `json:"redirects"`
				Rendered struct {
					Screenshot string `json:"screenshot"`
				} `json:"rendered"`
				Path   string `json:"path"`
				Status struct {
					Code int    `json:"code"`
				} `json:"status"`
			} `json:"response"`
		} `json:"data"`
	} `json:"result"`
}

// DataTableRow represents a row in the DataTable
type DataTableRow struct {
	URL        string
	Title      string
	Redirects  string
	Screenshot string
	Path       string
	StatusCode int
}

func readJSONData(filePath string) ([]*JSONData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data []*JSONData

	for {
		var jsonData JSONData
		err := decoder.Decode(&jsonData)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			log.Println("Error decoding JSON:", err)
			continue
		}

		data = append(data, &jsonData)
	}

	return data, nil
}

func formatRedirects(redirects []struct {
	Headers     struct {
		Location   string `json:"location"`
		Connection string `json:"connection"`
	} `json:"headers"`
	Status      struct {
		Code int `json:"code"`
	} `json:"status"`
	RedirectURI string `json:"redirect_uri"`
}) string {
	result := ""
	for _, redirect := range redirects {
		result += fmt.Sprintf("Location: %s\nConnection: %s\n\n", redirect.Headers.Location, redirect.Headers.Connection)
	}
	return result
}

func generateHTML(rows []DataTableRow) error {
	tmpl := template.Must(template.New("datatable").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
		    <meta name="description" content="">
		    <meta name="viewport" content="width=device-width, initial-scale=1">
		    <meta property="og:title" content="">
		    <meta property="og:type" content="">
		    <meta property="og:url" content="">
		    <meta property="og:image" content="">
		    <link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.10.25/css/jquery.dataTables.min.css">
			<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/buttons/2.2.3/css/buttons.dataTables.min.css">
			<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/css/all.min.css">
			<link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/bootstrap@4.6.1/dist/css/bootstrap.min.css">
			<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
			<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.6.1/dist/js/bootstrap.bundle.min.js"></script>
			<script src="https://cdn.datatables.net/1.10.25/js/jquery.dataTables.min.js"></script>
			<script src="https://cdn.datatables.net/buttons/2.2.3/js/dataTables.buttons.min.js"></script>
			<script src="https://cdn.datatables.net/buttons/2.2.3/js/buttons.html5.min.js"></script>
			<script src="https://cdn.datatables.net/buttons/2.2.3/js/buttons.print.min.js"></script>


			<script type="text/javascript" class="init">
				$(document).ready(function() {
				$('#datatable').DataTable( {
				dom: 'Bfrtip',
				lengthMenu: [[100, 250, 500, -1],[100, 250, 500, 'All']],
				buttons: [
							'copy', 'csv', 'excel', 'print'
						 ]
					} );
				} );
		</script>
		</head>
		<body>
			<table id="datatable" class="display" style="width:100%">
				<thead>
					<tr>
						<th>Response URL</th>
						<th>Title</th>
						<th>Redirects</th>
						<th>Screenshot</th>
						<th>Path</th>
						<th>Status Code</th>
					</tr>
				</thead>
				<tbody>
					{{range .}}
						<tr>
							<td><textarea rows="5" cols="30" readonly>{{.URL}}</textarea></td>
							<td><textarea rows="5" cols="30" readonly>{{.Title}}</textarea></td>
							<td><textarea rows="5" cols="30" readonly>{{.Redirects}}</textarea></td>
							<td><a href="{{.URL}}" target="_blank" rel="noreferrer"><img src="{{.Screenshot}}"></a></td>
							<td><textarea rows="5" cols="30" readonly>{{.Path}}</textarea></td>
							<td>{{.StatusCode}}</td>
						</tr>
					{{end}}
				</tbody>
			</table>
			<script>
				$(document).ready(function() {
					$('#datatable').DataTable();
				});
			</script>
		</body>
		</html>
	`))

	file, err := os.Create("output.html")
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, rows)
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	fmt.Println("HTML file generated successfully.")

	return nil
}

func main() {
	data, err := readJSONData("data.json")
	if err != nil {
		log.Fatal("Failed to read JSON data:", err)
	}

	var rows []DataTableRow

	for _, jsonData := range data {
		redirects := formatRedirects(jsonData.Result.Data.Response.Redirects)

		row := DataTableRow{
			URL:        jsonData.Result.Data.Response.URL,
			Title:      jsonData.Result.Data.Response.Title,
			Redirects:  redirects,
			Screenshot: jsonData.Result.Data.Response.Rendered.Screenshot,
			Path:       jsonData.Result.Data.Response.Path,
			StatusCode: jsonData.Result.Data.Response.Status.Code,
		}

		rows = append(rows, row)
	}

	err = generateHTML(rows)
	if err != nil {
		log.Fatal("Failed to generate HTML:", err)
	}
}
