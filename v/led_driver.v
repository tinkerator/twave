/*
 * Code placed in public domain.
 */

`timescale 10ns/100ps

module led_driver(
    input wire clk,
    input wire reset__disable,
    output reg reset__ack,
    output reg led);

    /*
     * Locally managed data:
     */
    reg [11:0] counter;
    reg [11:0] thresh;

    /*
     * Logic for this module:
     */
    always @(posedge clk)
      begin
        reset__ack <= #1 ~ reset__disable;
        if (! reset__disable)
          begin
            led <= #1 1'b0;
            thresh <= #1 12'H1;
            counter <= #1 12'H0;
          end
        else if (thresh == counter)
          begin
            led <= #1 ~ led;
            counter <= #1 12'H0;
            thresh <= #1 (thresh + 12'H3);
          end
        else
          begin
            counter <= #1 (counter + 12'H1);
          end
      end

endmodule /* led_driver */
