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
import io.ktor.client.features.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json


// RPC is the API client.
// url is the required API endpoint address.
class RPC(val url: String) {
    // AuthToken is an optional authentication token.
    var authToken: String? = null

    // client is used for making requests, defaulting to URLSession.shared.
    val client = HttpClient(CIO)

`
var end = `
    // call implementation.
    suspend fun call(
        method: String, input: String
    ): String {

        val url = this.url + "/" + method
        try {
            val r = this.client.post<HttpResponse>(url) {
                headers {
                    append(HttpHeaders.ContentType, "application/json")
                    if (authToken != null) {
                        append(HttpHeaders.Authorization, "Bearer ${authToken}")
                    }
                }
                body = input
            }
            return r.readText()
        } catch (e: ClientRequestException) {
            val body = e.response.readText()
            val json = Json { ignoreUnknownKeys = true }.decodeFromString<ResponseError>(body)
            throw RPCError(
                status = e.response.status.description,
                statusCode = e.response.status.value,
                type = json.type,
                msg = json.message
            )
        }
    }
}

data class RPCError(
    val status: String,
    val statusCode: Int,
    val type: String,
    val msg: String
) : Exception()

@Serializable
private data class ResponseError(val type: String, val message: String)
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
        val s = Json { ignoreUnknownKeys = true }.encodeToString(input)
        call(method = "%s", input = s)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, m.Name)

}

func writeOutputOnlyMethod(w io.Writer, m schema.Method) {
	camel := strcase.ToCamel(m.Name)
	lcamel := strcase.ToLowerCamel(m.Name)
	template := `    suspend fun %s(): %sOutput {
        val out = call(method = "%s", input = "")
        return Json { ignoreUnknownKeys = true }.decodeFromString(out)
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
        val s = Json { ignoreUnknownKeys = true }.encodeToString(input)
        val out = call(method = "%s", input = s)
        return Json { ignoreUnknownKeys = true }.decodeFromString(out)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, camel, m.Name)

}
