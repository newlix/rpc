import Foundation

// Client is the API client.
struct Client {
    // url is the required API endpoint address.
    let url: URL

    // AuthToken is an optional authentication token.
    let authToken: String?

    // session is the client used for making requests, defaulting to URLSession.shared.
    let session: URLSession = URLSession.shared

    // addItem adds an item to the list.
    func addItem(input: AddItemInput, complete: @escaping (_ error: Error?) -> ()) {
        func done(_ none: None?, _ err: Error?) {
            complete(err)
        }
        call(session: self.session, authToken: self.authToken, endpoint: self.url, method: "add_item", input: input, complete: done)
    }

    // getItems returns all items in the list.
    func getItems(complete: @escaping (_ output: GetItemsOutput?, _ err: Error?) -> ()) {
        call(session: self.session, authToken: self.authToken, endpoint: self.url, method: "get_items", input: None.only, complete: complete)
    }

    // removeItem removes an item from the to-do list.
    func removeItem(input: RemoveItemInput, complete: @escaping (_ output: RemoveItemOutput?, _ error: Error?) -> Void) {
        call(session: self.session, authToken: self.authToken, endpoint: self.url, method: "remove_item", input: input, complete: complete)
    }

}

let encoder: JSONEncoder = { () -> JSONEncoder in
    let encoder = JSONEncoder()
    encoder.keyEncodingStrategy = .convertToSnakeCase
    return encoder
}()

let decoder: JSONDecoder = { () -> JSONDecoder in
    let decoder = JSONDecoder()
    decoder.keyDecodingStrategy = .convertFromSnakeCase
    return decoder
}()

struct RPCError: Codable, Error {
    let status: String
    let statusCode: Int
    let type: String
    let message: String
}

struct None: Codable {
    static let only: None = None()
}

// call implementation.
private func call<Input, Output>(session: URLSession, authToken: String?, endpoint: URL, method: String, input: Input, complete: @escaping (_ output: Output?, _ error: Error?) -> Void) where Input: Codable, Output: Codable {

    var r = URLRequest(url: URL(string: method, relativeTo: endpoint)!)
    r.httpMethod = "POST"
    r.setValue("Application/json", forHTTPHeaderField: "Content-Type")
    if let token = authToken {
        r.setValue("Bearer " + token, forHTTPHeaderField: "Authorization")
    }

    if !(Input.self is None.Type) {
        do {
            r.httpBody = try encoder.encode(input)
        } catch {
            //todo
        }
    }

    session.dataTask(with: r, completionHandler: { (data, response, resError) in
        guard let data = data, let httpResponse = response as? HTTPURLResponse, resError == nil else {
            print("No valid response: endpoint = \(endpoint), method = \(method)")
            complete(nil, resError)
            return
        }

        // error
        let code = httpResponse.statusCode
        if code >= 300 {
            do {
                let err = try decoder.decode(RPCError.self, from: data)
                complete(nil, err)
            } catch {
                let status = HTTPURLResponse.localizedString(forStatusCode: code)
                let err = RPCError(status: status, statusCode: code, type: "", message: "")
                complete(nil, err)
            }
        }

        // output
        if Output.self is None.Type {
            complete(nil, nil)
        } else {
            do {
                let out = try decoder.decode(Output.self, from: data)
                complete(out, nil)
            } catch {
                complete(nil, error)
            }
        }
    })
}