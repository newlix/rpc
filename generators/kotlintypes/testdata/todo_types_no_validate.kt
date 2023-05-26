import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

/**
 * Item is a to-do item.
 * @property createdAt is the time the to-do item was created.
 * @property id is the id of the item. This field is read-only.
 * @property text is the to-do item text. This field is required.
 */
@Serializable
data class Item(
    @SerialName("created_at") var createdAt: String = "",
    @SerialName("id") val id: Int = 0,
    @SerialName("text") var text: String = ""
)

/**
 * addItem input params.
 * @property item is the item to add. This field is required.
 */
@Serializable
data class AddItemInput(
    @SerialName("item") var item: String = ""
)

/**
 * getItems output params.
 * @property items is the list of to-do items.
 */
@Serializable
data class GetItemsOutput(
    @SerialName("items") var items: Array<Item> = arrayOf()
)

/**
 * removeItem input params.
 * @property id is the id of the item to remove.
 */
@Serializable
data class RemoveItemInput(
    @SerialName("id") var id: Int = 0
)

/**
 * removeItem output params.
 * @property item is the item removed.
 */
@Serializable
data class RemoveItemOutput(
    @SerialName("item") var item: Item = Item()
)

