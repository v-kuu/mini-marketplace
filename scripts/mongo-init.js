db = db.getSiblingDB("mini-marketplace");

db.createCollection("products", {
	validator: {
		$jsonSchema: {
			bsonType: "object",
			required: ["_id", "name", "price"],
			properties: {
				_id: {
					bsonType: "string",
					description: "required string"
				},
				name: {
					bsonType: "string",
					description: "required string"
				},
				price: {
					bsonType: "long",
					description: "required integer"
				}
			}
		}
	}
});

db.products.createIndex({name: 1}, {unique: true})
