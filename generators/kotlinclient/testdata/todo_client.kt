
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

    // addItem adds an item to the list.
    suspend fun addItem(input: AddItemInput) {
        val s = Json.encodeToString(intput)
        call(endpoint = this.url, method = "add_item", input = s)
    }

    // getItems returns all items in the list.
    suspend fun getItems(): GetItemsOutput {
        val out = call(endpoint = this.url, method = "get_items", input = "")
        return Json.decodeFromString(out)
    }

    // removeItem removes an item from the to-do list.
    suspend fun removeItem(
        input: RemoveItemInput
    ): RemoveItemOutput {
        val s = Json.encodeToString(intput)
        val out = call(endpoint = this.url, method = "remove_item", input = s)
        return Json.decodeFromString(out)
    }


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
