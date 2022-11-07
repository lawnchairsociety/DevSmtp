namespace DevSmtp.Core.Commands
{
    public class RcptException : Exception
    {
        public RcptException(string message)
            : base(message)
        {
        }

        public RcptException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
