#include "funcdata.h"
#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-4
    get_tls(CX)
    MOVL    g(CX), AX
    MOVL    AX, ret+0(FP)
    RET
