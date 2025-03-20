select schema::Function {
  Json := <json>.name,
}
filter .name = 'std::str_repeat' or .name = 'std::str_lower'
order by .name
