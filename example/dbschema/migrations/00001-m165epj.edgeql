CREATE MIGRATION m165epj7s23o3bk4jgcbcevzaebvjspab3wc5ojsoq6nevtew2fmpq
    ONTO initial
{
  CREATE FUTURE simple_scoping;
  CREATE TYPE default::Movie {
      CREATE REQUIRED PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY title: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY year: std::int64;
  };
};
