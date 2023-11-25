package kotlinclient

import (
	"fmt"
	"io"

	"github.com/newlix/rpc/schema"
	"github.com/iancoleman/strcase"
)

var start = `
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody

data class RPCError(
    val status: String,
    val statusCode: Int,
    val type: String,
    val msg: String
) : Exception()

@Serializable
private data class ResponseError(val type: String, val message: String)


// RPC is the API client.
// url is the required API endpoint address.
class RPC(val endpoint: String) {
    val decoder = Json { ignoreUnknownKeys = true }

    // AuthToken is an optional authentication token.
    var authToken: String? = null

    // client is used for making requests, defaulting to URLSession.shared.
    val client = OkHttpClient()

    private suspend fun call(
        method: String, input: String
    ): String {
        return withContext(Dispatchers.IO) {
            val url = endpoint + "/" + method

            val request = Request.Builder()
                .url(url)
                .post(input.toRequestBody())
                .addHeader("Content-Type", "application/json")
            if (authToken != null) {
                request.addHeader("Authorization", "Bearer ${authToken}")
            }

            return@withContext client.newCall(request.build()).execute().use { response ->
                val body: String = response.body!!.string()
                if (!response.isSuccessful) {
                    val json = decoder.decodeFromString<ResponseError>(body)
                    throw RPCError(
                        status = response.message,
                        statusCode = response.code,
                        type = json.type,
                        msg = json.message
                    )
                }
                return@use body
            }
        }
    }
`

// Generate writes the Go client implementations to w.
func Generate(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf
	// w = os.Stderr
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

	out(w, "}\n")

	return nil
}

func writeInputOnlyMethod(w io.Writer, m schema.Method) {
	camel := strcase.ToCamel(m.Name)
	lcamel := strcase.ToLowerCamel(m.Name)
	template := `    suspend fun %s(input: %sInput) {
        val s = decoder.encodeToString(input)
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
        return decoder.decodeFromString(out)
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
        val s = decoder.encodeToString(input)
        val out = call(method = "%s", input = s)
        return decoder.decodeFromString(out)
    }
`
	fmt.Fprintf(w, template, lcamel, camel, camel, m.Name)

}
