const TOKEN = '114514:TELEGRAM-BOT_TOKEN';

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
});

async function handleRequest(request) {
  const update = await request.json();
  const message = update.message;
  const bot = new Telegram(TOKEN);
  let smms
  if (message) {
    smms = new Smms(await USERTOKEN.get(message.from.id));
  } else {
    smms = new Smms();
  }
  // 命令
  if (message !== undefined && message.text !== undefined) {
    const command = message.text.split(' ')
    switch (command[0]) {
      case '/get':
        const usertoken = await USERTOKEN.get(message.from.id)
        return await bot.sendMessage(message.from.id, '`' + usertoken + '`', 'markdown', true, true, message.message_id)
      case '/set':
        if (command.length == 1) {
          return await bot.sendMessage(message.from.id, 'Format: /set xxxxxx', 'markdown', true, true, message.message_id)
        }
        if (command[1].length == 32) {
          await USERTOKEN.put(message.from.id, command[1])
          return await bot.sendMessage(message.from.id, 'API token saved.', 'markdown', true, true, message.message_id)
        }
        return await bot.sendMessage(message.from.id, 'API token length invalid.', 'markdown', true, true, message.message_id)
      case '/del':
        await USERTOKEN.put(message.from.id, command[1])
        return await bot.sendMessage(message.from.id, 'API token deleted.', 'markdown', true, true, message.message_id)
      default:
        return await bot.sendMessage(message.from.id, 'Unknown command', 'markdown', true, true, message.message_id)
    }
  }
  // 文件形式的图片
  if (message !== undefined && message.document !== undefined) {
    if (!message.document.mime_type.startsWith('image/')) {
      return await bot.sendMessage(message.from.id, 'File has an invalid extension.', 'markdown', true, true, message.message_id)
    }
    const fileID = message.document.file_id
    const filePath = await bot.getFile(fileID)
    const filePathJSON = await filePath.json()
    const url = `https://api.telegram.org/file/bot${TOKEN}/${filePathJSON.result.file_path}`
    const image = await fetch(url)
    const result = await smms.upload(await image.arrayBuffer())
    const resultJSON = await result.json()
    if (resultJSON.success) {
      return await bot.sendMessage(message.from.id, '`' + resultJSON.data.url + '`', 'markdown', true, true, message.message_id, {
        inline_keyboard: [[{
          text: resultJSON.data.hash,
          callback_data: resultJSON.data.hash
        }]]
      })
    }
    return await bot.sendMessage(message.from.id, '`' + resultJSON.message + '`', 'markdown', true, true, message.message_id)
  }
  // 图片
  if (message !== undefined && message.photo !== undefined) {
    const fileID = message.photo[message.photo.length - 1].file_id
    const filePath = await bot.getFile(fileID)
    const filePathJSON = await filePath.json()
    const url = `https://api.telegram.org/file/bot${TOKEN}/${filePathJSON.result.file_path}`
    const image = await fetch(url)
    const result = await smms.upload(await image.arrayBuffer())
    const resultJSON = await result.json()
    if (resultJSON.success) {
      return await bot.sendMessage(message.from.id, '`' + resultJSON.data.url + '`', 'markdown', false, false, message.message_id, {
        inline_keyboard: [[{
          text: resultJSON.data.hash,
          callback_data: resultJSON.data.hash
        }]]
      })
    }
    return await bot.sendMessage(message.from.id, '`' + resultJSON.message + '`', 'markdown', false, false, message.message_id)
  }
  // Callback 删除图片
  if (update.callback_query !== undefined) {
    await smms.delete(update.callback_query.data)
    return await bot.editMessageText(update.callback_query.message.chat.id, update.callback_query.message.message_id, null, "Photo Deleted!")
  }
  return new Response()
}

class Telegram {
  constructor(token) {
    this.api = 'https://api.telegram.org/bot' + token;
  }
  async sendMessage(chat_id, text, parse_mode, disable_web_page_preview, disable_notification,
    reply_to_message_id, reply_markup) {
    return await fetch(this.api + '/sendMessage', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        chat_id: chat_id,
        text: text,
        parse_mode: parse_mode,
        disable_web_page_preview: disable_web_page_preview,
        disable_notification: disable_notification,
        reply_to_message_id: reply_to_message_id,
        reply_markup: reply_markup
      })
    })
  }
  async editMessageText(chat_id, message_id, inline_message_id, text, parse_mode, disable_web_page_preview,
    reply_markup) {
    return await fetch(this.api + '/editMessageText', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        chat_id: chat_id,
        message_id: message_id,
        inline_message_id: inline_message_id,
        text: text,
        parse_mode: parse_mode,
        disable_web_page_preview: disable_web_page_preview,
        reply_markup: reply_markup
      })
    })
  }
  async getFile(file_id) {
    return await fetch(this.api + '/getFile', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        file_id: file_id
      })
    })
  }
}

class Smms {
  constructor(token) {
    this.token = token;
    this.api = 'https://sm.ms/api/v2'
  }
  async upload(image) {
    const f = fd(image);
    const headers = new Headers({
      'Content-Type': 'multipart/form-data; boundary=-114514',
      'User-Agent': 'sm_ms_bot/cloudflareworkers (https://github.com/imlonghao/smms-bot)'
    });
    if (this.token !== null) {
      headers.append('Authorization', this.token)
    }
    return await fetch(this.api + '/upload/', {
      method: 'POST',
      headers: headers,
      body: f
    })
  }
  async delete(hash) {
    const headers = new Headers({
      'User-Agent': 'sm_ms_bot/cloudflareworkers (https://github.com/imlonghao/smms-bot)'
    });
    return await fetch(this.api + '/delete/' + hash, {
      headers: headers
    })
  }
}

var _appendBuffer = function (buffer1, buffer2) {
  var tmp = new Uint8Array(buffer1.byteLength + buffer2.byteLength);
  tmp.set(new Uint8Array(buffer1), 0);
  tmp.set(new Uint8Array(buffer2), buffer1.byteLength);
  return tmp.buffer;
};

function fd(image) {
  const enc = new TextEncoder();
  let txt = '---114514\r\n';
  txt += 'Content-Disposition: form-data; name="smfile"; filename="uploaded_by_sm_ms_bot"\r\n';
  txt += '\r\n';
  const p1 = enc.encode(txt);
  txt = '\r\n';
  txt += '---114514--\r\n';
  const p2 = enc.encode(txt);
  return _appendBuffer(p1, _appendBuffer(image, p2))
}
