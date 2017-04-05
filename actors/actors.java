
import java.io.IOException;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLConnection;
import java.util.Arrays;
import java.util.Map;
import java.util.TreeMap;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import javax.swing.JFrame;
import javax.swing.JOptionPane;
import javax.swing.JPanel;

//Garret Stevens

public class Main {
    
    static Map<String, BlockingQueue<Object[]> > queues = new TreeMap<>();
    
    static Theora.Context ctx = new Theora.Context();
    static void theoraActor(BlockingQueue<Object[]> Q){
        //does the theora decoding
        //remember, decoding takes place wherever you call Theora.getFrame()
        
        byte[] byte_ = new byte[BUFFER_SIZE];
        
        try{
            while(true){
                byte_[0] = (byte) Q.take()[0]; //locking
                //System.out.print((char) byte_[0]);
                ctx.enqueueData(byte_, 0, BUFFER_SIZE);
                //send message to color actor
                //System.out.println(ctx.getBufferedDataSize());
                
            }
        }catch(InterruptedException e){
            return;
        }
    }
    
    static void colorActor(BlockingQueue<Object[]> Q){
        //color space conversion
        Theora.ImagePlane[] ip;
        try{
            Object[] o = Q.take(); //locking
            while(true){
                //do stuff
                ip = ctx.getFrame(); //locking
            }
        }catch(InterruptedException e){
            return;
        }
    }
    
    static void drawActor(BlockingQueue<Object[]> Q){
        //does ui stuff. draws the frames and even ends the program with "THE END"
        try{
            while(true){
                Object[] o = Q.take(); //locking
                //do stuff
            }
        }catch(InterruptedException e){
            return;
        }
    }
    
    static final int BUFFER_SIZE = 1;
    static void netActor(BlockingQueue<Object []> Q){
        //does the network downloading
        byte[] buffer = new byte[BUFFER_SIZE];
        URL u = null;
        InputStream i = null;
        while(true){
            try{
                u = new URL(x);
            }catch(MalformedURLException e){
                System.out.println(e);
                System.exit(0);
            }
            
            try{
                URLConnection c = u.openConnection();
                i = c.getInputStream();
                                
                //print out the data, yo
                while(i.read(buffer) > 0){
                    for(int j=0;j<BUFFER_SIZE;j++){
                        //System.out.print((char) buffer[j]);
                        
                        try{ //send byte to theora actor
                            queues.get("TA").put(new Object [] {buffer[j]});
                        }catch(InterruptedException e){
                            System.out.println(e);
                            System.exit(0);
                        }
                        
                        //if(buffer[j] == '\0')
                            //break;
                    }
                }
                
            }catch(IOException e){
                System.out.println(e);
                System.exit(0);
            }
            break;
        }
    }
    
    //Globals
    static JFrame w;
    //static byte[] buffer = new byte[BUFFER_SIZE];
    //String x = JOptionPane.showInputDialog("Enter a website:"); //prompt for a url
    static String x = "http://selenium.ssucet.org/04.OGG";
    
    public static void main(String[] args){
        
        try{
            queues.put("TA",new LinkedBlockingQueue<>());
            queues.put("CA",new LinkedBlockingQueue<>());
            queues.put("DA",new LinkedBlockingQueue<>());
            queues.put("NA",new LinkedBlockingQueue<>());
            Thread t1 = new Thread( () -> { theoraActor(queues.get("TA")); } );
            Thread t2 = new Thread( () -> { colorActor(queues.get("CA")); } );
            Thread t3 = new Thread( () -> { drawActor(queues.get("DA")); } );
            Thread t4 = new Thread( () -> { netActor(queues.get("NA")); } );
            t1.start();
            t2.start();
            t3.start();
            t4.start();
            
            //sample site: http://selenium.ssucet.org
            //String x = JOptionPane.showInputDialog("Enter a website:"); //prompt for a url
            //String x = "http://selenium.ssucet.org/04.OGG";
            //URL u = null;
            //InputStream i = null;
            
            //starts the thing
            queues.get("NA").put(new Object [] {x});
            queues.get("CA").put(new Object[] {"go"});
            
            
            w = new JFrame(x);
            w.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
            JPanel p = new JPanel();
            w.add(p);
            w.setSize(1280, 720);
            w.setVisible(true);
            
            
            
            //Theora.Context ctx = new Theora.Context();
            
            //System.out.println(Arrays.toString(buffer));
            //ctx.enqueueData(buffer, 0, BUFFER_SIZE);
            //Theora.ImagePlane[] ip = ctx.getFrame();
            
            
            //queues.get("A1").put(new Object[]{"GO"});
            t1.join();
            t2.join();
            t3.join();
            t4.join();
        }catch(InterruptedException e){
            System.exit(0);
        }
        
        /*
            Create a video-display program that uses actors (and message queues
            to communicate between the actors). Since this problem's dataflow is
            fairly linear, it will look a lot like a pipeline paradigm.

            You should have separate actors for the network downloading, for the
            Theora decoding (remember, decoding takes place wherever you call
            Theora.getFrame()), for the color space conversion, and for the UI.

            When the video is over, the screen should display "The End". You
            should look up how to draw text to a Graphics object in the Javadocs. ;-)
        */
    }
    
}
