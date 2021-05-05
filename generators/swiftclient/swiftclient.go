package swiftclient

import (
	"fmt"
	"io"

	"github.com/apex/rpc/schema"
	"github.com/iancoleman/strcase"
)

var start = `import Foundation

// %s is the API client.
struct %s {
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

`
var end = `
// call implementation.
private func call<Input, Output>(endpoint: URL, method: String, input: Input, complete: @escaping (_ output: Output?, _ error: Error?) -> Void) where Input: Codable, Output: Codable {

    var url = endpoint
    url.appendPathComponent(method, isDirectory: false)

    r.httpMethod = "POST"
    r.setValue("application/json", forHTTPHeaderField: "Content-Type")
    if let token = self.authToken {
        r.setValue("Bearer " + token, forHTTPHeaderField: "Authorization")
    }

    do {
        r.httpBody = try self.encoder.encode(input)
    } catch {
        complete(nil, error)
    }

    self.session.dataTask(with: r) { (data, response, resError) in
        let response: HTTPURLResponse! = response as? HTTPURLResponse
        if response == nil {
            complete(nil, "not http response: respone: \(response) err:(\(resError)")
        }


        // error
        let code = httpResponse.statusCode
        let status = HTTPURLResponse.localizedString(forStatusCode: code)
        if code >= 300 {
            do {
                let body = try self.decoder.decode(ResponseErrorBody.self, from: data)
                let err = HTTPError(status: status, statusCode: code, type: body.type, message: body.message)
                complete(nil, err)
            } catch {
                complete(nil, error)
            }
        }

        // output
        do {
            let out = try self.decoder.decode(Output.self, from: data)
            complete(out, nil)
        } catch {
            complete(nil, error)
        }
    }
}
}

struct HTTPError: Error {
let status: String
let statusCode: Int
let type: String
let message: String
}

struct ResponseErrorBody: Codable {
let type: String
let message: String
}

extension String: Error {

}
`

// Generate writes the Go client implementations to w.
func Generate(w io.Writer, s *schema.Schema, client string) error {
	out := fmt.Fprintf

	out(w, start, client, client)

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
	template := `    func %s(input: %sInput, complete: @escaping (_ error: Error?) -> ()) {
        callInputOnly(endpoint: self.url, method: "%s", input: input, complete: complete)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, m.Name)

}

func writeOutputOnlyMethod(w io.Writer, m schema.Method) {
	camel := strcase.ToCamel(m.Name)
	lcamel := strcase.ToLowerCamel(m.Name)
	template := `    func %s(complete: @escaping (_ output: %sOutput?, _ err: Error?) -> ()) {
        callOutputOnly(endpoint: self.url, method: "%s", complete: complete)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, m.Name)

}

func writeMethod(w io.Writer, m schema.Method) {
	camel := strcase.ToCamel(m.Name)
	lcamel := strcase.ToLowerCamel(m.Name)
	template := `    func %s(input: %sInput, complete: @escaping (_ output: %sOutput?, _ error: Error?) -> Void) {
        call(endpoint: self.url, method: "%s", input: input, complete: complete)
    }

`
	fmt.Fprintf(w, template, lcamel, camel, camel, m.Name)

}
