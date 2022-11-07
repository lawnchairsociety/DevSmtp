namespace DevSmtp.Core.Commands
{
    public class RsetException : Exception
    {
        public RsetException(string message)
            : base(message)
        {
        }

        public RsetException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
