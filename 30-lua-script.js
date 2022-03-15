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

    console.log("DOC : https://redis.io/topics/lua-api");

    let script1 = `
      local value1 = '-' .. ARGV[1] .. '-'
      redis.call('set', KEYS[1], ARGV[1]);
      redis.call('set', KEYS[2], value1);
      local a = {
          redis.call('get', KEYS[1]), 
          redis.call('get', KEYS[2])
        }
      return a[2]
      `;

    let r1 = await client.sendCommand([
      "EVAL",
      script1,
      "2",
      "key1",
      "key2",
      "titi",
    ]);
    console.log(r1);

    let id = await client.sendCommand(["SCRIPT", "LOAD", script1]);

    console.log("script id =", id);
    let r1bis;

    r1bis = await client.sendCommand([
      "EVALSHA",
      id,
      "2",
      "key1",
      "key2",
      "titi",
    ]);
    console.log("r1bis =", r1bis);

    r1bis = await client.sendCommand([
      "EVALSHA",
      id,
      "2",
      "key1",
      "key2",
      "toto",
    ]);
    console.log("r1bis =", r1bis);

    let r2 = await client.sendCommand([
      "EVAL",
      `
      local a = ARGV[1];
      return a .. ' toto'
      `,
      "0",
      "hello",
    ]);
    console.log(r2);

    client.quit();
  } catch (error) {
    console.error("error=", error);
    process.exit();
  }
})();
