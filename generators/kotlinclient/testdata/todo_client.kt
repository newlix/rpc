
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

    // addItem adds an item to the list.
    suspend fun addItem(input: AddItemInput) {
        val s = Json { ignoreUnknownKeys = true }.encodeToString(input)
        call(method = "add_item", input = s)
    }

    // getItems returns all items in the list.
    suspend fun getItems(): GetItemsOutput {
        val out = call(method = "get_items", input = "")
        return Json { ignoreUnknownKeys = true }.decodeFromString(out)
    }

    // removeItem removes an item from the to-do list.
    suspend fun removeItem(
        input: RemoveItemInput
    ): RemoveItemOutput {
        val s = Json { ignoreUnknownKeys = true }.encodeToString(input)
        val out = call(method = "remove_item", input = s)
        return Json { ignoreUnknownKeys = true }.decodeFromString(out)
    }


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
