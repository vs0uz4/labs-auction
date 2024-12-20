db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE || "auction_db");

db.createCollection("users");
db.createCollection("auctions");
db.createCollection("bids");

db.auctions.createIndex({ status: 1, category: 1 });
db.auctions.createIndex({ product_name: "text" });

db.bids.createIndex({ auction_id: 1, amount: -1 });

db.users.insertMany([
    {
      _id: "d290f1ee-6c54-4b01-90e6-d701748f0851",
      name: "John Doe",
    },
    {
      _id: "93fb1e9c-523f-4d92-80b4-0f7ba12fef56",
      name: "Jane Smith",
    },
    {
      _id: "4be43d3d-5f47-4881-a07b-8b5d3c5296c1",
      name: "Alice Johnson",
    },
]);