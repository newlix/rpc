import Foundation

// Item is a to-do item.
struct Item: Codable {
    // createdAt is the time the to-do item was created.
    var createdAt: Date

    // id is the id of the item. This field is read-only.
    var id: Int

    // text is the to-do item text. This field is required.
    var text: String

    enum CodingKeys: String, CodingKey {
        case createdAt = "created_at"
        case id = "id"
        case text = "text"
    }
}

// AddItemInput params.
struct AddItemInput: Codable {
    // item is the item to add. This field is required.
    var item: String

    enum CodingKeys: String, CodingKey {
        case item = "item"
    }
}

// GetItemsOutput params.
struct GetItemsOutput: Codable {
    // items is the list of to-do items.
    var items: [Item]

    enum CodingKeys: String, CodingKey {
        case items = "items"
    }
}

// RemoveItemInput params.
struct RemoveItemInput: Codable {
    // id is the id of the item to remove.
    var id: Int

    enum CodingKeys: String, CodingKey {
        case id = "id"
    }
}

// RemoveItemOutput params.
struct RemoveItemOutput: Codable {
    // item is the item removed.
    var item: Item

    enum CodingKeys: String, CodingKey {
        case item = "item"
    }
}

