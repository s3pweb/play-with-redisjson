let createClient = require("redis").createClient;

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

    let now = new Date();

    console.log(
      "WRITE ->",
      await client.json.set("input:5e986128435ad3fafbeb88dc", "$", {
        resourceId: "5e98612a9a72a30010ec03f4", // could be a secondary index
        entityId: "5e9572de2b4aae0010433600", // could be a secondary index

        description: "BLABLABLA",
        speed: { ts: now.getTime(), value: 25 },
        course: { ts: now.getTime(), value: 200 },
        sensor: [1001, 1002, 1003],
      })
    );

    console.log(
      "1 ->" +
        JSON.stringify(await client.json.get("input:5e986128435ad3fafbeb88dc"))
    );

    console.log(
      "2 ->" +
        JSON.stringify(
          await client.json.get("input:5e986128435ad3fafbeb88dc", {
            path: ".description",
          })
        )
    );

    console.log(
      "3 ->" +
        JSON.stringify(
          await client.json.get("input:5e986128435ad3fafbeb88dc", {
            path: ".course.ts",
          })
        )
    );

    client.quit();
  } catch (error) {
    console.error(error);
  }
})();
