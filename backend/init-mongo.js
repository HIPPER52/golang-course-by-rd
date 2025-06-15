db = db.getSiblingDB('support')

db.operators.insertOne({
    id: "0001M41752KX6RE8QFA3D3QW32",
    username: "admin",
    email: "admin@admin.com",
    pwd_hash: "$2a$10$bU0KdhIYlnquaLJd.iqqueGYwC.glSYNz8QQZK.hoDUTh8dchy2Si", // password: 12345678
    role: "admin",
    created_at: new Date()
})
