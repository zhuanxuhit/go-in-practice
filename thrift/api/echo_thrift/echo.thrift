namespace go echo

struct EchoReq {
    1: string name
}
struct EchoRes {
    1: string msg
}

service Echo {
    EchoRes echo(1: EchoReq req);
}
