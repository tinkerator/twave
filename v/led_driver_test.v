/*
 * Code placed in public domain.
 */

`timescale 10ns/100ps

module led_driver_test();

    /*
     * Locally managed data:
     */
    wire clk;
    wire led;
    reg reset__disable;
    wire reset__ack;
    reg [11:0] script;

    /*
     * Inner modules referenced by this module:
     *
     * => led_driver
     */

    led_driver target(
        clk,

        reset__disable,
        reset__ack,

        led);


    /*
     * Logic for this module:
     */
    assign clk = script[1];

    initial
      begin
        $dumpvars();
        reset__disable <= #1 1'b1;
        script <= #1 12'H0;
      end

    always @(posedge clk, script)
      begin
        script <= #1 (script + 12'H1);
        case (script)
        12'H6:
          begin
            reset__disable <= #1 1'b0;
          end
        12'H14:
          begin
            if (! reset__ack)
              begin
                $display("led.bb(51):RESULT=FAIL:1 @ testing failed reset ack");
                $finish_and_return(1);
              end
            reset__disable <= #1 1'b1;
          end
        12'H12c:
          begin
            $display("led.bb(56):RESULT=PASS:0 @ testing completed");
            $finish_and_return(0);
          end
        endcase
      end


endmodule /* led_driver_test */
