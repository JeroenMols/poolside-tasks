package routes

const userNameRegex = `^[a-zA-Z0-9 ]{5,32}$`
const accountNumberRegex = `^usr_[23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy=]{22}$`

// TODO: allow more characters
const todoDescriptionRegex = `^[a-zA-Z0-9 ]{1,256}$`

// Fake data
const fakeToken = "tkn_aaaaaaaaaaaaaaaaaaaaaa"
const fakeWrongToken = "tkn_bbbbbbbbbbbbbbbbbbbbbb"

const fakeUserId = "usr_aaaaaaaaaaaaaaaaaaaaaa"
const fakeWrongUserId = "usr_bbbbbbbbbbbbbbbbbbbbbb"
const fakeTodoListId = "lst_aaaaaaaaaaaaaaaaaaaaaa"
const fakeWrongTodoListId = "lst_bbbbbbbbbbbbbbbbbbbbbb"
const fakeTodoListId2 = "lst_cccccccccccccccccccccc"
const fakeTodoId = "tdo_aaaaaaaaaaaaaaaaaaaaaa"
const fakeWrongTodoId = "tdo_bbbbbbbbbbbbbbbbbbbbbb"
