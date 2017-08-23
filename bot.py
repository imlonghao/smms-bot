#!/usr/bin/env python3
#
# Copyright (c) 2017 imlonghao <shield@fastmail.com>
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

import logging
from requests_futures.sessions import FuturesSession
from telegram import InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import Updater, Filters, MessageHandler, CallbackQueryHandler
from telegram.ext.dispatcher import run_async
from os import mkdir, remove, environ

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.DEBUG)
logger = logging.getLogger(__name__)


def download(bot, file_id):
    bot.getFile(file_id).download('/tmp/smms-bot/%s' % file_id)


def upload(file_id):
    f = {
        'smfile': open('/tmp/smms-bot/%s' % file_id, 'rb')
    }
    r = requests.post('https://sm.ms/api/upload', files=f)
    return r.result().json()


@run_async
def error_handler(bot, update, error):
    logger.exception(error)


@run_async
def upload_handler(bot, update):
    try:
        file_id = update.message.document.file_id
        if not update.message.document.mime_type.startswith('image/'):
            return update.message.reply_text('File has an invalid extension.', quote=True)
    except:
        file_id = update.message.photo[-1].file_id
    download(bot, file_id)
    uploader = upload(file_id)
    remove('/tmp/smms-bot/%s' % file_id)
    if uploader['code'] == 'error':
        update.message.reply_text(uploader['msg'], quote=True)
    else:
        kb = [[InlineKeyboardButton('Click Here To Delete', callback_data=uploader['data']['hash'])]]
        update.message.reply_text('`%s`' % uploader['data']['url'], quote=True, parse_mode='markdown',
                                  disable_web_page_preview=True, reply_markup=InlineKeyboardMarkup(kb))


@run_async
def callback_handler(bot, update):
    key = update.callback_query.data
    requests.get('https://sm.ms/api/delete/%s' % key)
    update.callback_query.message.edit_text('Photo Deleted!')


if __name__ == '__main__':
    try:
        mkdir('/tmp/smms-bot')
    except FileExistsError:
        pass
    requests = FuturesSession(max_workers=10)
    updater = Updater(environ['TG_TOKEN'], workers=10)
    dp = updater.dispatcher
    dp.add_handler(MessageHandler(Filters.document, upload_handler))
    dp.add_handler(MessageHandler(Filters.photo, upload_handler))
    dp.add_handler(CallbackQueryHandler(callback_handler))
    dp.add_error_handler(error_handler)
    updater.start_polling()
    updater.idle()
