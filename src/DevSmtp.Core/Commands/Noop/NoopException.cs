namespace DevSmtp.Core.Commands
{
    public class NoopException : Exception
    {
        public NoopException(string message)
            : base(message)
        {
        }

        public NoopException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
