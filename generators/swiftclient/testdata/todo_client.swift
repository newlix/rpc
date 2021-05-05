import Foundation

// Client is the API client.
struct Client {
    // encoder is the conventional json encoder
    let encoder = JSONEncoder()

    // decoder is the conventional json decoder
    let decoder = JSONDecoder()

    // url is the required API endpoint address.
    let url: URL

    // AuthToken is an optional authentication token.
    var authToken: String?

    // session is the client used for making requests, defaulting to URLSession.shared.
    let session: URLSession = URLSession.shared

    // addItem adds an item to the list.
    func addItem(input: AddItemInput, complete: @escaping (_ error: Error?) -> ()) {
        call(endpoint: self.url, method: "add_item", input: input, complete: { (_: Nothing?, err: Error?) in complete(err) })
    }

    // getItems returns all items in the list.
    func getItems(complete: @escaping (_ output: GetItemsOutput?, _ err: Error?) -> ()) {
        call(endpoint: self.url, method: "get_items", input: Nothing(), complete: complete)
    }

    // removeItem removes an item from the to-do list.
    func removeItem(input: RemoveItemInput, complete: @escaping (_ output: RemoveItemOutput?, _ error: Error?) -> Void) {
        call(endpoint: self.url, method: "remove_item", input: input, complete: complete)
    }


    // call implementation.
    private func call<Input, Output>(endpoint: URL, method: String, input: Input, complete: @escaping (_ output: Output?, _ error: Error?) -> Void) where Input: Codable, Output: Codable {

        var url = endpoint
        url.appendPathComponent(method, isDirectory: false)

        var r = URLRequest(url: url)
        r.httpMethod = "POST"
        r.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let token = self.authToken {
            r.setValue("Bearer " + token, forHTTPHeaderField: "Authorization")
        }

        do {
            if !(input is Nothing) {
                r.httpBody = try self.encoder.encode(input)
            }
        } catch {
            complete(nil, error)
        }

        self.session.dataTask(with: r) { (data, response, resError) in
            let response: HTTPURLResponse! = response as? HTTPURLResponse
            if response == nil {
                complete(nil, "not http response: respone: \(String(describing: response)) err:(\(String(describing: resError))")
                return
            }


            // error
            let code = response.statusCode
            let status = HTTPURLResponse.localizedString(forStatusCode: code)
            if code >= 300 {
                do {
                    let body = try self.decoder.decode(ResponseErrorBody.self, from: data ?? Data())
                    let err = HTTPError(status: status, statusCode: code, type: body.type, message: body.message)
                    complete(nil, err)
                } catch {
                    complete(nil, error)
                }
                return
            }

            // output
            do {
                if Output.self is Nothing.Type {
                    complete(nil, nil)
                } else {
                    let out = try self.decoder.decode(Output.self, from: data ?? Data())
                    complete(out, nil)
                }
            } catch {
                complete(nil, error)
            }
        }.resume()
    }
}

struct HTTPError: Error {
    let status: String
    let statusCode: Int
    let type: String
    let message: String
}

struct ResponseErrorBody: Codable {
    let type: String
    let message: String
}

extension String: Error {

}

struct Nothing: Codable {

}
