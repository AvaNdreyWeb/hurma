db.createUser({
    user: "admin",
    pwd: "password",
    roles: [
      {
        role: "readWrite",
        db: "hurma"
      }
    ]
  });
  
  db.getSiblingDB("hurma").createCollection("users");
  db.getSiblingDB("hurma").createCollection("links");