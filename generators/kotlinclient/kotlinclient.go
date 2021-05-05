package kotlinclient

import (
	"fmt"
	"io"

	"github.com/apex/rpc/schema"
	"github.com/iancoleman/strcase"
)

var start = `
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json


// Client is the API client.
// url is the required API endpoint address.
class Client(val url: String) {
    // AuthToken is an optional authentication token.
    var authToken: String? = null

    // client is used for making requests, defaulting to URLSession.shared.
    var client = HttpClient(CIO)

`
var end = `
    // call implementation.
    suspend fun call(
        endpoint: String, method: String, input: String
    ): String {

        var url = this.url + "/" + endpoint

        val r = this.client.post<HttpResponse>(url) {
            append(HttpHeaders.ContentType, "application/json")
            if (this.authToken != null) {
                append(HttpHeaders.Authorization, "Bearer ${this.authToken}")
            }
            body = input
        }

        if (r.status.value >= 300) {
            var err = Json.decodeFromString<ResponseError>(r.readText())
            throw HTTPExcepion(
                status = r.status.description,
                statusCode = r.status.value,
                type = err.type,
                msg = err.message
            )
        }
        return r.readText()
    }
}

data class HTTPExcepion(
    val status: String,
    val statusCode: Int,
    val type: String,
    val msg: String
) : Exception()

@Serializable
data class ResponseError(val type: String, val message: String)
`

// Generate writes the Go client implementations to w.
func Generate(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf

	out(w, start)

	for _, m := range s.Methods {
		name := strcase.ToLowerCamel(m.Name)
		out(w, "    // %s %s\n", name, m.Description)

		if len(m.Inputs) > 0 && len(m.Outputs) == 0 {
			writeInputOnlyMethod(w, m)
			continue
		}

		if len(m.Inputs) == 0 && len(m.Outputs) > 0 {
			writeOutputOnlyMethod(w, m)
			continue
		}

		if len(m.Inputs) > 0 && len(m.Outputs) > 0 {
			writeMethod(w, m)
			continue
		}
	}

	out(w, end)

	return nil
}

func writeInputOnlyMethod(w io.Writer, m schema.Method) {
	camel := strcase.ToCamel(m.Name)
	lcamel := strcase.ToLowerCamel(m.Name)
	template := `    suspend fun %s(input: %sInput) {
        val s = Json.encodeToString(intput)
        call(endpoint = this.url, method = "%s", input = s)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, m.Name)

}

func writeOutputOnlyMethod(w io.Writer, m schema.Method) {
	camel := strcase.ToCamel(m.Name)
	lcamel := strcase.ToLowerCamel(m.Name)
	template := `    suspend fun %s(): %sOutput {
        val out = call(endpoint = this.url, method = "%s", input = "")
        return Json.decodeFromString(out)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, m.Name)

}

func writeMethod(w io.Writer, m schema.Method) {
	camel := strcase.ToCamel(m.Name)
	lcamel := strcase.ToLowerCamel(m.Name)
	template := `    suspend fun %s(
        input: %sInput
    ): %sOutput {
        val s = Json.encodeToString(intput)
        val out = call(endpoint = this.url, method = "%s", input = s)
        return Json.decodeFromString(out)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, camel, m.Name)

}
