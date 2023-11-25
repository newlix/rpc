
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
    // addItem adds an item to the list.
    suspend fun addItem(input: AddItemInput) {
        val s = decoder.encodeToString(input)
        call(method = "add_item", input = s)
    }

    // getItems returns all items in the list.
    suspend fun getItems(): GetItemsOutput {
        val out = call(method = "get_items", input = "")
        return decoder.decodeFromString(out)
    }

    // removeItem removes an item from the to-do list.
    suspend fun removeItem(
        input: RemoveItemInput
    ): RemoveItemOutput {
        val s = decoder.encodeToString(input)
        val out = call(method = "remove_item", input = s)
        return decoder.decodeFromString(out)
    }
}
