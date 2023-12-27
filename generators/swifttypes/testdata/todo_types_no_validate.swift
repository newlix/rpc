import Foundation

// Item is a to-do item.
struct Item: Codable {
    // createdAt is the time the to-do item was created.
    var createdAt: Date = Date()

    // id is the id of the item. This field is read-only.
    var id: Int = 0

    // text is the to-do item text. This field is required.
    var text: String = ""

    enum CodingKeys: String, CodingKey {
        case createdAt = "created_at"
        case id = "id"
        case text = "text"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let createdAt = try container.decodeIfPresent(String.self, forKey: .createdAt) {
            self.createdAt = createdAt
        }
        if let id = try container.decodeIfPresent(String.self, forKey: .id) {
            self.id = id
        }
        if let text = try container.decodeIfPresent(String.self, forKey: .text) {
            self.text = text
        }
    }
}

// AddItemInput params.
struct AddItemInput: Codable {
    // item is the item to add. This field is required.
    var item: String = ""

    enum CodingKeys: String, CodingKey {
        case item = "item"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let item = try container.decodeIfPresent(String.self, forKey: .item) {
            self.item = item
        }
    }
}

// GetItemsOutput params.
struct GetItemsOutput: Codable {
    // items is the list of to-do items.
    var items: [Item] = []

    enum CodingKeys: String, CodingKey {
        case items = "items"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let items = try container.decodeIfPresent(String.self, forKey: .items) {
            self.items = items
        }
    }
}

// RemoveItemInput params.
struct RemoveItemInput: Codable {
    // id is the id of the item to remove.
    var id: Int = 0

    enum CodingKeys: String, CodingKey {
        case id = "id"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let id = try container.decodeIfPresent(String.self, forKey: .id) {
            self.id = id
        }
    }
}

// RemoveItemOutput params.
struct RemoveItemOutput: Codable {
    // item is the item removed.
    var item: Item = Item()

    enum CodingKeys: String, CodingKey {
        case item = "item"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let item = try container.decodeIfPresent(String.self, forKey: .item) {
            self.item = item
        }
    }
}

