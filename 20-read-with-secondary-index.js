let createClient = require("redis").createClient;
let SchemaFieldTypes = require("redis").SchemaFieldTypes;

(async () => {
  try {
    const client = createClient({
      url: "redis://localhost:6379",
      socket: {
        tls: false,
      },
    });

    client.on("error", (err) => console.log("Redis Client Error", err));

    await client.connect();

    console.log("connected...");

    let indexes = await client.ft._list();
    let indexName = "idx:input";
    if (indexes.indexOf(indexName) >= 0) {
      console.log("Drop index ->", await client.ft.dropIndex(indexName));
    }

    console.log(
      "Create index ->",
      await client.ft.create(
        "idx:input",
        {
          "$.resourceId": {
            type: SchemaFieldTypes.TAG,
            AS: "resourceId",
            SEPARATOR: ";",
          },
          "$.entityId": {
            type: SchemaFieldTypes.TAG,
            AS: "entityId",
            SEPARATOR: ";",
          },
          "$.speed.value": {
            type: SchemaFieldTypes.NUMERIC,
            AS: "speed",
          },
          "$.course.value": {
            type: SchemaFieldTypes.NUMERIC,
            AS: "course",
          },
        },
        {
          ON: "JSON",
          PREFIX: "input",
        }
      )
    );

    let now = new Date();

    console.log(
      "WRITE ->",
      await client.json.set("input:5e98612a9a72a30010ec0001", "$", {
        id: "5e98612a9a72a30010ec0001",
        resourceId: "5e98612a9a72a30010ec03f5", // could be a secondary index
        entityId: "5e9572de2b4aae0010433600", // could be a secondary index

        speed: {
          ts: now.getTime(),
          value: 50,
        },
        course: {
          ts: now.getTime(),
          value: 100,
        },
      })
    );

    now = new Date();

    console.log(
      "WRITE ->",
      await client.json.set("input:5e98612a9a72a30010ec0003", "$", {
        id: "5e98612a9a72a30010ec0003",
        resourceId: "5e98612a9a72a30010ec03ba", // could be a secondary index
        entityId: "5e9572de2b4aae0010433600", // could be a secondary index

        speed: {
          ts: now.getTime(),
          value: 75,
        },
        course: {
          ts: now.getTime(),
          value: 20,
        },
      })
    );

    now = new Date();

    console.log(
      "WRITE ->",
      await client.json.set("input:5e98612a9a72a30010ec0004", "$", {
        id: "5e98612a9a72a30010ec0004",
        resourceId: "5e98612a9a72a30010ec03ba", // could be a secondary index
        entityId: "5e9572de2b4aae0010433600", // could be a secondary index

        speed: {
          ts: now.getTime(),
          value: 18,
        },
        course: {
          ts: now.getTime(),
          value: 45,
        },
      })
    );

    console.log(
      "READ ONE ->" +
        JSON.stringify(
          await client.ft.search(
            "idx:input",
            "@resourceId:{5e98612a9a72a30010ec03f5}"
          )
        )
    );

    console.log(
      "READ MULTIPLE ->" +
        JSON.stringify(
          await client.ft.search(
            "idx:input",
            "@resourceId:{5e98612a9a72a30010ec03ba | 5e98612a9a72a30010ec03f5}"
          )
        )
    );

    console.log(
      "READ BETWEEN VALUES 1 ->" +
        JSON.stringify(
          await client.ft.search(
            "idx:input",
            "@resourceId:{5e98612a9a72a30010ec03f4} @speed:[0 70]"
          )
        )
    );

    console.log(
      "READ BETWEEN VALUES 2 ->" +
        JSON.stringify(
          await client.ft.search("idx:input", "@speed:[0 70] @course:[0 60]")
        )
    );

    client.quit();
  } catch (error) {
    console.error(error);
  }
})();
